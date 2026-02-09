document.addEventListener("DOMContentLoaded", async () => {
    const params = new URLSearchParams(window.location.search);
    const id = params.get("id");

    if (!id) {
        document.getElementById("product-container").innerHTML = "<p>Product not found.</p>";
        return;
    }

    try {
        const res = await fetch(`/products/${id}`);
        if (!res.ok) throw new Error("Product not found");
        
        const p = await res.json();
        
        const container = document.getElementById("product-container");
        const imgUrl = p.image_url || "https://via.placeholder.com/400";

        container.innerHTML = `
            <div class="product-image">
                <img src="${imgUrl}" alt="${p.name}">
            </div>
            <div class="product-info">
                <h1>${p.name}</h1>
                <p class="price">${p.price_kzt} KZT</p>
                <div class="meta">
                    <span>Category: ${p.category}</span> | 
                    <span>Size: ${p.size}</span> | 
                    <span>Stock: ${p.stock_quantity}</span>
                </div>
                <p class="description">
                   This is a high-quality ${p.name} made for your comfort. 
                   Perfect for any occasion.
                </p>
                <button class="btn-large" onclick="addToCart(${p.id})">Add to Cart</button>
            </div>
        `;
    } catch (err) {
        console.error(err);
        document.getElementById("product-container").innerHTML = "<p>Error loading product.</p>";
    }
});

function addToCart(id) {
    let cart = JSON.parse(localStorage.getItem("cart") || "[]");
    
    cart.push(id);
    
    localStorage.setItem("cart", JSON.stringify(cart));
    
    alert("Item added to cart!");
}