document.addEventListener("DOMContentLoaded", async () => {
    const token = localStorage.getItem("token");
    const authSection = document.getElementById("auth-section");
    const profileSection = document.getElementById("profile-section");

    if (!token) {
        if (authSection) authSection.style.display = "flex";
        setupAuthForms();
    } else {
        if (profileSection) profileSection.style.display = "block";
        loadUserProfile(token);
        loadOrders(token);
    }
});

function setupAuthForms() {
    const loginForm = document.getElementById("login-form");
    if (loginForm) {
        loginForm.addEventListener("submit", async (e) => {
            e.preventDefault();
            const email = document.getElementById("login-email").value;
            const password = document.getElementById("login-password").value;
            await handleAuth("/login", { email, password });
        });
    }

    const signupForm = document.getElementById("signup-form");
    if (signupForm) {
        signupForm.addEventListener("submit", async (e) => {
            e.preventDefault();
            const userData = {
                name: document.getElementById("s-name").value,
                email: document.getElementById("s-email").value,
                password: document.getElementById("s-password").value,
                address: document.getElementById("s-address").value
            };
            const success = await handleAuth("/signup", userData);
            if(success) alert("Account created! Please login.");
        });
    }
}

async function handleAuth(endpoint, data) {
    try {
        const res = await fetch(endpoint, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data)
        });

        if (res.ok) {
            if (endpoint === "/login") {
                const resData = await res.json();
                localStorage.setItem("token", resData.token);
                if(resData.role) localStorage.setItem("role", resData.role);
                window.location.reload();
            }
            return true;
        } else {
            alert("Action failed. Check details.");
            return false;
        }
    } catch (err) {
        console.error(err);
        return false;
    }
}

async function loadUserProfile(token) {
    try {
        const res = await fetch("/api/users/me", {
            headers: { "Authorization": "Bearer " + token }
        });
        if (res.ok) {
            const user = await res.json();
            const nameElem = document.getElementById("user-name");
            const emailElem = document.getElementById("user-email");
            const addrElem = document.getElementById("user-address");
            
            if (nameElem) nameElem.innerText = `Hello, ${user.name}`;
            if (emailElem) emailElem.innerText = user.email;
            if (addrElem) addrElem.innerText = user.address;
        } else {
            logout();
        }
    } catch (err) { console.error(err); }
}

async function loadOrders(token) {
    const list = document.getElementById("order-list");
    if (!list) return;
    
    try {
        const res = await fetch("/api/orders", {
            headers: { "Authorization": "Bearer " + token }
        });

        if (res.ok) {
            const orders = await res.json();
            
            if (!orders || orders.length === 0) {
                list.innerHTML = "<li>No orders found.</li>";
                return;
            }

            list.innerHTML = ""; 

            orders.forEach(order => {
                const date = new Date(order.created_at).toLocaleDateString();
                let statusColor = 'green';
                if (order.status === 'Pending') statusColor = 'orange';
                if (order.status === 'Cancelled') statusColor = 'red';

                const li = document.createElement("li");
                li.style.display = "flex";
                li.style.justifyContent = "space-between";
                li.style.padding = "10px";
                li.style.borderBottom = "1px solid #eee";
                
                li.innerHTML = `
                    <div>
                        <strong>Order #${order.id}</strong>
                        <br><small style="color: #888">${date}</small>
                    </div>
                    <div style="text-align: right; margin-right: 15px;">
                        <span style="font-weight: bold;">${order.total_price} KZT</span>
                        <br><small style="color: ${statusColor}; font-weight: bold;">${order.status}</small>
                    </div>
                    <div>
                        <button class="btn-view" onclick="viewOrder(${order.id}, '${order.status}')">View</button>
                    </div>
                `;
                list.appendChild(li);
            });
        }
    } catch (err) {
        console.error(err);
        list.innerHTML = "<li>Error loading orders.</li>";
    }
}

async function viewOrder(orderId, status) {
    const token = localStorage.getItem("token");
    const modal = document.getElementById("order-modal");
    const tbody = document.getElementById("modal-items");
    const cancelBtn = document.getElementById("cancel-order-btn");
    const title = document.getElementById("modal-title");

    if (title) title.innerText = `Order #${orderId} Details`;
    if (tbody) tbody.innerHTML = "<tr><td colspan='4'>Loading...</td></tr>";
    if (modal) modal.style.display = "block";

    if (cancelBtn) {
        if (status === "Pending") {
            cancelBtn.style.display = "block";
            cancelBtn.onclick = function() { cancelOrder(orderId); };
        } else {
            cancelBtn.style.display = "none";
        }
    }

    try {
        const res = await fetch(`/api/orders/${orderId}`, {
            headers: { "Authorization": "Bearer " + token }
        });

        if (res.ok && tbody) {
            const items = await res.json();
            tbody.innerHTML = "";
            items.forEach(item => {
                tbody.innerHTML += `
                    <tr>
                        <td>${item.product_name}</td>
                        <td>${item.quantity}</td>
                        <td>${item.price} KZT</td>
                        <td>${item.price * item.quantity} KZT</td>
                    </tr>
                `;
            });
        } else if (tbody) {
            tbody.innerHTML = "<tr><td colspan='4'>Failed to load items.</td></tr>";
        }
    } catch (err) { console.error(err); }
}

async function cancelOrder(orderId) {
    if (!confirm("Are you sure you want to cancel this order?")) return;

    const token = localStorage.getItem("token");
    try {
        const res = await fetch(`/api/orders/${orderId}/cancel`, {
            method: "PUT",
            headers: { "Authorization": "Bearer " + token }
        });

        if (res.ok) {
            alert("Order Cancelled!");
            closeModal();
            loadOrders(token); 
        } else {
            alert("Failed to cancel order.");
        }
    } catch (err) { console.error(err); }
}

function closeModal() {
    const modal = document.getElementById("order-modal");
    if (modal) modal.style.display = "none";
}

function logout() {
    localStorage.removeItem("token");
    localStorage.removeItem("role");
    window.location.reload();
}