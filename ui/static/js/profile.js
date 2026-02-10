document.addEventListener("DOMContentLoaded", async () => {
    const token = localStorage.getItem("token");
    const authSection = document.getElementById("auth-section");
    const profileSection = document.getElementById("profile-section");

    if (!token) {
        authSection.style.display = "flex";
        setupAuthForms();
    } else {
        profileSection.style.display = "block";
        
        loadUserProfile(token);
        
        loadOrders(token);
    }
});

function setupAuthForms() {
    document.getElementById("login-form").addEventListener("submit", async (e) => {
        e.preventDefault();
        const email = document.getElementById("login-email").value;
        const password = document.getElementById("login-password").value;
        await handleAuth("/login", { email, password });
    });

    document.getElementById("signup-form").addEventListener("submit", async (e) => {
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
                if (resData.role) localStorage.setItem("role", resData.role);
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
            document.getElementById("user-name").innerText = `Hello, ${user.name}`;
            document.getElementById("user-email").innerText = user.email;
            document.getElementById("user-address").innerText = user.address;
        } else {
            logout();
        }
    } catch (err) {
        console.error(err);
    }
}

async function loadOrders(token) {
    const list = document.getElementById("order-list");
    
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
                const statusColor = order.status === 'Pending' ? 'orange' : 'green';

                const li = document.createElement("li");
                
                li.style.padding = "10px";
                li.style.borderBottom = "1px solid #eee";
                li.style.display = "flex";
                li.style.justifyContent = "space-between";

                li.innerHTML = `
                    <div>
                        <strong>Order #${order.id}</strong>
                        <br><small style="color: #888">${date}</small>
                    </div>
                    <div style="text-align: right;">
                        <span style="font-weight: bold;">${order.total_price} KZT</span>
                        <br><small style="color: ${statusColor}">${order.status}</small>
                    </div>
                `;
                list.appendChild(li);
            });
        }
    } catch (err) {
        console.error("Error loading orders:", err);
        list.innerHTML = "<li>Error loading orders.</li>";
    }
}

function logout() {
    localStorage.removeItem("token");
    localStorage.removeItem("role");
    window.location.reload();
}