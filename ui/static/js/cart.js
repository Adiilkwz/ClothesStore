document.addEventListener("DOMContentLoaded", async () => {
    const cartIds = JSON.parse(localStorage.getItem("cart") || "[]");
    
    if (cartIds.length === 0) {
        const emptyMsg = document.getElementById("empty-msg");
        if (emptyMsg) emptyMsg.style.display = "block";
        return;
    }

    const counts = {};
    cartIds.forEach(id => { counts[id] = (counts[id] || 0) + 1; });

    const uniqueIds = Object.keys(counts);
    let totalPrice = 0;
    const tbody = document.getElementById("cart-table-body");

    try {
        const productPromises = uniqueIds.map(id => 
            fetch(`/products/${id}`).then(r => {
                if (!r.ok) throw new Error(`Product ${id} not found`);
                return r.json();
            })
        );
        const products = await Promise.all(productPromises);

        if (tbody) {
            tbody.innerHTML = ""; 
            products.forEach(p => {
                const qty = counts[p.id];
                const lineTotal = p.price_kzt * qty;
                totalPrice += lineTotal;

                tbody.innerHTML += `
                    <tr>
                        <td>
                            <img src="${p.image_url || 'https://via.placeholder.com/50'}" alt="${p.name}" style="width:50px; margin-right:10px;">
                            ${p.name}
                        </td>
                        <td>${p.price_kzt} KZT</td>
                        <td>${qty}</td>
                        <td>${lineTotal} KZT</td>
                        <td><button class="remove-btn" onclick="removeItem(${p.id})">X</button></td>
                    </tr>
                `;
            });
        }

        const totalItemsElem = document.getElementById("total-items");
        const totalPriceElem = document.getElementById("total-price");
        if (totalItemsElem) totalItemsElem.innerText = cartIds.length;
        if (totalPriceElem) totalPriceElem.innerText = totalPrice + " KZT";

    } catch (err) {
        console.error("Error loading cart:", err);
    }

    const checkoutBtn = document.getElementById("checkout-btn");
    if (checkoutBtn) {
        checkoutBtn.addEventListener("click", async () => {
            const token = localStorage.getItem("token");
            if (!token) {
                alert("Please login to checkout!");
                window.location.href = "/user/login"; 
                return;
            }

            const items = uniqueIds.map(id => ({
                product_id: parseInt(id),
                quantity: counts[id]
            }));

            try {
                const res = await fetch('/api/orders', {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": "Bearer " + token
                    },
                    body: JSON.stringify({ items: items })
                });

                if (res.ok) {
                    alert("Order Placed Successfully!");
                    localStorage.removeItem("cart");
                    window.location.href = "/user/profile";
                } else {
                    const errorData = await res.text();
                    alert("Order failed: " + errorData);
                }
            } catch (err) {
                console.error("Checkout error:", err);
                alert("Connection error. Please try again.");
            }
        });
    }
});

function removeItem(idToRemove) {
    let cart = JSON.parse(localStorage.getItem("cart") || "[]");

    const newCart = cart.filter(id => id !== idToRemove);
    localStorage.setItem("cart", JSON.stringify(newCart));
    window.location.reload();
}
