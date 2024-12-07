// Wait for the DOM content to load before executing the script
document.addEventListener('DOMContentLoaded', async function () {
    const userId = localStorage.getItem('user_id'); // Retrieve the user ID from local storage
    if (!userId) {
        alert("Not logged in"); // Notify if the user is not logged in
        window.location.href = 'index.html'; // Redirect to the login page
        return;
    }
    
    try {
        // Fetch user data from the server
        const response = await fetch(`http://localhost:8080/user/${userId}`);
        if (response.ok) {
            const user = await response.json(); // Parse the user data
            // Populate profile fields with user data
            document.getElementById('email').value = user.email || 'N/A';
            document.getElementById('phone').value = user.phone || 'N/A';
            document.getElementById('membership_tier').value = user.membership_tier || 'N/A';
        } else {
            alert('Failed to fetch user data!'); // Notify if fetching user data fails
        }

        // Open the profile edit popup
        document.getElementById('editProfile').addEventListener('click', () => {
            document.getElementById('editProfilePopup').style.display = "flex";
        });

        // Close the profile edit popup
        document.querySelector('.close-btn').addEventListener('click', () => {
            document.getElementById('editProfilePopup').style.display = "none";
        });

        // Handle profile update form submission
        document.getElementById('editProfileForm').addEventListener('submit', async (e) => {
            e.preventDefault(); // Prevent default form submission

            // Collect updated profile data
            const newEmail = document.getElementById('newEmail').value;
            const newPhone = document.getElementById('newPhone').value;
            const newPassword = document.getElementById('newPassword').value;
            
            try {
                // Send an update request to the server
                const response = await fetch(`http://localhost:8080/user/${userId}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        email: newEmail,
                        phone: newPhone,
                        password: newPassword,
                    }),
                });

                if (response.ok) {
                    alert('Profile updated successfully!'); // Notify user on success
                    document.getElementById('email').value = newEmail;
                    document.getElementById('phone').value = newPhone;
                    document.getElementById('membership_tier').value = user.membership_tier; 
                    document.getElementById('editProfilePopup').style.display = "none"; // Close the popup
                } else {
                    alert('Failed to update profile!'); // Notify if update fails
                }
            } catch (error) {
                console.error('Error updating profile:', error); // Log errors
            }
        });

        // Load and display rental history
        document.getElementById('viewHistory').addEventListener('click', async () => {
            const rentalHistorySection = document.getElementById('rentalHistorySection');
            const rentalHistoryList = rentalHistorySection.querySelector('#rentalHistoryList');
            rentalHistoryList.innerHTML = ''; // Clear previous history content

            try {
                // Fetch rental history data
                const historyResponse = await fetch(`http://localhost:8080/rental-history/${userId}`);
                if (historyResponse.ok) {
                    const rentalHistory = await historyResponse.json();

                    if (rentalHistory.length > 0) {
                        rentalHistory.forEach(item => {
                            // Create rental history entry
                            const listItem = document.createElement('div');
                            listItem.innerHTML = `
                                <p><strong>Vehicle:</strong> ${item.carPlate || 'Unknown'}</p>
                                <p><strong>Start Time:</strong> ${item.startTime ? new Date(item.startTime).toLocaleString() : 'N/A'}</p>
                                <p><strong>End Time:</strong> ${item.endTime ? new Date(item.endTime).toLocaleString() : 'N/A'}</p>
                                <p><strong>Total Price:</strong> $${item.totalPrice || '0.00'}</p>
                                <p><strong>Status:</strong> ${item.status || 'Unknown'}</p>
                            `;

                            // Add action buttons for active reservations
                            if (item.status === 'Active') {
                                const buttonsContainer = document.createElement('div');
                                buttonsContainer.innerHTML = `
                                    <button class="edit-btn" data-id="${item.id}">Edit Reservation</button>
                                    <button class="complete-btn" data-id="${item.id}">Complete Reservation</button>
                                `;
                                listItem.appendChild(buttonsContainer);
                            }

                            rentalHistoryList.appendChild(listItem);
                        });

                        // Add event listeners for reservation actions (edit, complete)
                        document.querySelectorAll('.edit-btn').forEach(button => {
                            button.addEventListener('click', async (e) => {
                                const reservationId = e.target.dataset.id;
                                // Fetch and display reservation details for editing
                                // ...
                            });
                        });

                        document.querySelectorAll('.complete-btn').forEach(button => {
                            button.addEventListener('click', async (e) => {
                                const reservationId = e.target.dataset.id;
                                // Complete reservation and initiate payment
                                // ...
                            });
                        });
                    } else {
                        rentalHistoryList.innerHTML = '<p>No rental history available.</p>'; // Notify if history is empty
                    }
                } else {
                    alert('Failed to fetch rental history!'); // Notify if fetching history fails
                }
            } catch (error) {
                console.error('Error fetching rental history:', error); // Log errors
            }
        });

    } catch (error) {
        console.error('Error:', error); // Log general errors
        alert('An error occurred while loading the page.');
    }
});
