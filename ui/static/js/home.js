document.addEventListener("DOMContentLoaded", async () => {
    try {
        const res = await fetch('/products');
        
        if (!res.ok) throw new Error("Failed to fetch");

        const products = await res.json();
        const grid = document.getElementById("product-grid");

        if (!products || products.length === 0) {
            grid.innerHTML = "<p>No products found. Admin needs to add some!</p>";
            return;
        }

        grid.innerHTML = "";

        products.forEach(p => {
            const img = p.image_url || "https://via.placeholder.com/200?text=No+Image";

            grid.innerHTML += `
                <div class="card">
                    <img src="${img}" alt="${p.name}" onerror="this.src='https://via.placeholder.com/200'">
                    <h3>${p.name}</h3>
                    <p>${p.price_kzt} KZT</p>
                    <a href="/product?id=${p.id}" class="btn">View Details</a>
                </div>
            `;
        });

    } catch (err) {
        console.error("Error:", err);
        document.getElementById("product-grid").innerHTML = 
            "<p style='color:red'>Error loading products. Is the server running?</p>";
    }
});