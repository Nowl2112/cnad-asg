document.addEventListener("DOMContentLoaded", () => {
    const loginForm = document.getElementById("loginForm");
    loginForm.addEventListener("submit", async (e) => {
        e.preventDefault();

        const email = document.getElementById("email").value;
        const password = document.getElementById("password").value;

        try {
            const response = await fetch("http://localhost:8080/users/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ email, password }),
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || "Login failed");
            }

            const data = await response.json();

            // Store the user ID and email in local storage
            localStorage.setItem("user_id", data.user_id);

            alert("Login successful!");
            window.location.href = "profile.html"; // Redirect to homepage
        } catch (err) {
            console.error("Error:", err);
            alert(err.message);
        }
    });
});
