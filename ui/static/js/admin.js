document.addEventListener("DOMContentLoaded", () => {
    fetchProducts();

    document.getElementById("add-product-form").addEventListener("submit", async (e) => {
        e.preventDefault();
        
        const token = localStorage.getItem("token");
        if (!token) {
            alert("You are not logged in!");
            return;
        }

        const product = {
            name: document.getElementById("p-name").value,
            price_kzt: parseInt(document.getElementById("p-price").value),
            image_url: document.getElementById("p-image").value,
            category: document.getElementById("p-category").value,
            size: document.getElementById("p-size").value,
            stock_quantity: parseInt(document.getElementById("p-stock").value)
        };

        try {
            const res = await fetch('/api/admin/products', {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": "Bearer " + token
                },
                body: JSON.stringify(product)
            });

            if (res.ok) {
                alert("Product Added!");
                document.getElementById("add-product-form").reset();
                fetchProducts();
            } else {
                alert("Failed. Check if you have Admin role.");
            }
        } catch (err) {
            console.error(err);
        }
    });
});

async function fetchProducts() {
    try {
        const res = await fetch('/products');
        const products = await res.json();
        
        const tbody = document.getElementById("admin-products-list");
        tbody.innerHTML = "";

        products.forEach(p => {
            tbody.innerHTML += `
                <tr>
                    <td>${p.id}</td>
                    <td>${p.name}</td>
                    <td>${p.price_kzt}</td>
                    <td>${p.stock_quantity}</td>
                    <td><button onclick="deleteProduct(${p.id})">Delete</button></td>
                </tr>
            `;
        });
    } catch (err) {
        console.error("Error loading admin products:", err);
    }
}

async function deleteProduct(id) {
    if(!confirm("Are you sure?")) return;
    
    const token = localStorage.getItem("token");
    await fetch(`/api/admin/products/${id}`, {
        method: "DELETE",
        headers: { "Authorization": "Bearer " + token }
    });
    fetchProducts();
}