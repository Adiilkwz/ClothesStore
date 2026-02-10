document.addEventListener("DOMContentLoaded", async () => {
    const params = new URLSearchParams(window.location.search);
    const id = params.get("id");

    if (!id) return;

    try {
        const res = await fetch(`/products/${id}`);
        if (!res.ok) throw new Error("Product not found");
        const p = await res.json();
        
        const container = document.getElementById("product-container");
        const imgUrl = p.image_url || "https://via.placeholder.com/400";

        let sizeOptions = '<option value="">Select Size</option>';
        if (p.size) {
            const sizes = p.size.split(",").map(s => s.trim());
            sizes.forEach(s => {
                sizeOptions += `<option value="${s}">${s}</option>`;
            });
        } else {
             sizeOptions = '<option value="OneSize">One Size</option>';
        }

        container.innerHTML = `
            <div class="product-image">
                <img src="${imgUrl}" alt="${p.name}">
            </div>
            <div class="product-info">
                <h1>${p.name}</h1>
                <p class="price">${p.price_kzt} KZT</p>
                
                <div class="meta">
                    <span>Category: ${p.category}</span> | 
                    <span>Stock: ${p.stock_quantity}</span>
                </div>

                <p class="description">
                   High-quality ${p.name}. Perfect for your collection.
                </p>

                <div class="selection-area">
                    <div class="selection-row">
                        <label for="size-select">Choose Size:</label>
                        <select id="size-select">
                            ${sizeOptions}
                        </select>
                    </div>

                    <div class="selection-row">
                        <label for="qty-input">Quantity:</label>
                        <input type="number" id="qty-input" value="1" min="1" max="${p.stock_quantity}">
                    </div>
                </div>

                <button class="btn-large" onclick="addToCart(${p.id}, '${p.name}', ${p.price_kzt})">Add to Cart</button>
            </div>
        `;
    } catch (err) {
        console.error(err);
        document.getElementById("product-container").innerHTML = "<p>Error loading product.</p>";
    }
});

function addToCart(id, name, price) {
    const sizeSelect = document.getElementById("size-select");
    const qtyInput = document.getElementById("qty-input");

    const selectedSize = sizeSelect.value;
    const quantity = parseInt(qtyInput.value);

    if (!selectedSize) {
        alert("Please select a size first!");
        return;
    }
    if (quantity < 1) {
        alert("Quantity must be at least 1");
        return;
    }

    const cartItem = {
        productId: id,
        name: name,
        price: price,
        size: selectedSize,
        quantity: quantity
    };

    let cart = JSON.parse(localStorage.getItem("cart_items") || "[]");

    const existingIndex = cart.findIndex(item => item.productId === id && item.size === selectedSize);

    if (existingIndex > -1) {
        cart[existingIndex].quantity += quantity;
    } else {
        cart.push(cartItem);
    }
    
    localStorage.setItem("cart_items", JSON.stringify(cart));
    
    alert(`Added ${quantity} x ${name} (${selectedSize}) to cart!`);
}