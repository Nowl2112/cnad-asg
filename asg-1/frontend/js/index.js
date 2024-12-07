// Run the code after the DOM content has fully loaded
document.addEventListener("DOMContentLoaded", () => {
    // Select the login form by its ID
    const loginForm = document.getElementById("loginForm");

    // Add an event listener to handle form submission
    loginForm.addEventListener("submit", async (e) => {
        // Prevent the default form submission behavior
        e.preventDefault();

        // Get the input values for email and password
        const email = document.getElementById("email").value;
        const password = document.getElementById("password").value;

        try {
            // Send a POST request to the login endpoint with the entered credentials
            const response = await fetch("http://localhost:8080/users/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json", 
                },
                body: JSON.stringify({ email, password }), 
            });

            // Check if the response indicates an error
            if (!response.ok) {
                const error = await response.json(); // Get the error message from the server
                throw new Error(error.message || "Login failed"); // Throw an error with the message
            }

            // Parse the response data if the request is successful
            const data = await response.json();

            // Store the user ID in localStorage for later use
            localStorage.setItem("user_id", data.user_id);

            // Notify the user of successful login and redirect to homepage
            alert("Login successful!");
            window.location.href = "homepage.html"; 
        } catch (err) {
            // Log and display any errors that occur
            console.error("Error:", err);
            alert(err.message);
        }
    });
});
