document.addEventListener("DOMContentLoaded", () => {
    const cart = JSON.parse(localStorage.getItem("cart_items") || "[]");
    
    const tbody = document.getElementById("cart-table-body");
    const emptyMsg = document.getElementById("empty-msg");
    const summaryBox = document.querySelector(".cart-summary");

    if (cart.length === 0) {
        emptyMsg.style.display = "block";
        document.querySelector(".cart-table").style.display = "none";
        
        if (summaryBox) summaryBox.style.display = "none";
        
        return;
    }

    let totalPrice = 0;
    let totalItems = 0;

    tbody.innerHTML = "";

    cart.forEach((item, index) => {
        const lineTotal = item.price * item.quantity;
        totalPrice += lineTotal;
        totalItems += item.quantity;

        tbody.innerHTML += `
            <tr>
                <td>
                    <strong>${item.name}</strong><br>
                    <small>Size: ${item.size}</small>
                </td>
                <td>${item.price} KZT</td>
                <td>${item.quantity}</td>
                <td>${lineTotal} KZT</td>
                <td><button class="remove-btn" onclick="removeItem(${index})">X</button></td>
            </tr>
        `;
    });

    document.getElementById("total-items").innerText = totalItems;
    document.getElementById("total-price").innerText = totalPrice + " KZT";

    document.getElementById("checkout-btn").addEventListener("click", async () => {
        const token = localStorage.getItem("token");
        if (!token) {
            alert("Please login to checkout!");
            window.location.href = "/profile";
            return;
        }

        const orderItems = cart.map(item => ({
            product_id: item.productId,
            quantity: item.quantity
        }));

        try {
            const res = await fetch('/api/orders', {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": "Bearer " + token
                },
                body: JSON.stringify({ items: orderItems })
            });

            if (res.ok) {
                const data = await res.json(); 
                
                alert(`Order Placed Successfully! (Order ID: ${data.order_id})\nCheck your email for confirmation.`);
                
                localStorage.removeItem("cart_items");
                window.location.href = "/profile";
            } else {
                alert("Order failed. Please try again.");
            }
        } catch (err) {
            console.error(err);
        }
    });
});

function removeItem(index) {
    let cart = JSON.parse(localStorage.getItem("cart_items") || "[]");
    cart.splice(index, 1);
    localStorage.setItem("cart_items", JSON.stringify(cart));
    window.location.reload();
}