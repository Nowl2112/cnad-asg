// Add an event listener to the registration form
document.getElementById("registerForm").addEventListener("submit", async (e) => {
    // Prevent the default form submission behavior to handle it via JavaScript
    e.preventDefault();

    // Retrieve values from input fields
    const email = document.getElementById("email").value; // User's email
    const password = document.getElementById("password").value; // User's password
    const phone = document.getElementById("phone").value; // User's phone number
    const membershipTier = "Basic"; // Default membership tier for new users

    try {
        // Make a POST request to the server to register the user
        const response = await fetch("http://localhost:8080/users/register", {
            method: "POST", // HTTP method
            headers: {
                "Content-Type": "application/json", // Content type is JSON
            },
            body: JSON.stringify({
                email, // Include email in the request body
                password, // Include password in the request body
                phone, // Include phone in the request body
                membership_tier: membershipTier, // Assign a default membership tier
            }),
        });

        // Check if the response status is not OK
        if (!response.ok) {
            // If registration fails, parse the error message from the server
            const error = await response.json();
            throw new Error(error.message || "Registration failed"); // Throw an error with the message
        }

        // If successful, parse the response data
        const data = await response.json();
        console.log("Registration successful:", data); // Log the success response for debugging

        // Inform the user of successful registration
        alert("Registration successful!");
    } catch (err) {
        // Handle any errors that occur during the registration process
        console.error("Error:", err); // Log the error for debugging
        alert(err.message); // Display the error message to the user
    }
});
