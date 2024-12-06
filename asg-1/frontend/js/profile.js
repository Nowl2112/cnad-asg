document.addEventListener('DOMContentLoaded', async function () {
    const userId = localStorage.getItem('user_id');
    if (!userId) {
        alert("Not logged in");
        window.location.href = 'index.html'; // Redirect to login if not logged in
        return;
    }

    try {
        // Fetch user data
        const response = await fetch(`http://localhost:8080/user/${userId}`);
        if (response.ok) {
            const user = await response.json();
            document.getElementById('email').value = user.email || 'N/A';
            document.getElementById('phone').value = user.phone || 'N/A';
            document.getElementById('membership_tier').value = user.membership_tier || 'N/A';
        } else {
            alert('Failed to fetch user data!');
        }

        // Event listener for editing profile
        document.getElementById('editProfile').addEventListener('click', () => {
            document.getElementById('editProfilePopup').style.display = "flex";
        });

        // Event listener for closing the edit profile popup
        document.querySelector('.close-btn').addEventListener('click', () => {
            document.getElementById('editProfilePopup').style.display = "none";
        });

        // Event listener for saving changes in the edit profile form
        document.getElementById('editProfileForm').addEventListener('submit', async (e) => {
            e.preventDefault();

            const newEmail = document.getElementById('newEmail').value;
            const newPhone = document.getElementById('newPhone').value;
            const newPassword = document.getElementById('newPassword').value;

            // Send the updated data to the server (assuming an endpoint for updating profile)
            try {
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
                    alert('Profile updated successfully!');
                    // Update the displayed profile information
                    document.getElementById('email').value = newEmail;
                    document.getElementById('phone').value = newPhone;
                    document.getElementById('membership_tier').value = user.membership_tier; // Keep membership tier as it is
                    document.getElementById('editProfilePopup').style.display = "none";
                } else {
                    alert('Failed to update profile!');
                }
            } catch (error) {
                console.error('Error updating profile:', error);
                alert('An error occurred while updating your profile.');
            }
        });

        // Event listener for rental history
// Event listener for rental history
document.getElementById('viewHistory').addEventListener('click', async () => {
    const rentalHistorySection = document.getElementById('rentalHistorySection');
    const rentalHistoryList = rentalHistorySection.querySelector('#rentalHistoryList');

    rentalHistoryList.innerHTML = ''; // Clear previous items

    try {
        const historyResponse = await fetch(`http://localhost:8080/rental-history/${userId}`);
        if (historyResponse.ok) {
            const rentalHistory = await historyResponse.json();
            if (rentalHistory.length > 0) {
                rentalHistory.forEach(item => {
                    const listItem = document.createElement('div');
                    listItem.innerHTML = `
                        <p><strong>Vehicle:</strong> ${item.carPlate || 'Unknown'}</p>
                        <p><strong>Start Time:</strong> ${item.startTime ? new Date(item.startTime).toLocaleString() : 'N/A'}</p>
                        <p><strong>End Time:</strong> ${item.endTime ? new Date(item.endTime).toLocaleString() : 'N/A'}</p>
                        <p><strong>Total Price:</strong> $${item.totalPrice || '0.00'}</p>
                        <p><strong>Status:</strong> ${item.status || 'Unknown'}</p>
                    `;

                    // Add Edit and Complete buttons if the reservation is active
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

                // Add event listeners to the buttons after rendering
                document.querySelectorAll('.edit-btn').forEach(button => {
                    button.addEventListener('click', async (e) => {
                        const reservationId = e.target.dataset.id;
                        // Handle the edit reservation logic (e.g., show an edit form or open a popup)
                        alert(`Edit reservation with ID: ${reservationId}`);
                    });
                });

                document.querySelectorAll('.complete-btn').forEach(button => {
                    button.addEventListener('click', async (e) => {
                        const reservationId = e.target.dataset.id;
                        // Call the API to complete the reservation
                        try {
                            const response = await fetch(`http://localhost:8082/reservations/${reservationId}/complete`, {
                                method: 'PUT', 
                            });
                            if (response.ok) {
                                alert(`Reservation ${reservationId} completed successfully!`);
                                // Optionally, refresh the rental history list or update the status locally
                            } else {
                                alert('Failed to complete reservation!');
                            }
                        } catch (error) {
                            console.error('Error completing reservation:', error);
                            alert('An error occurred while completing the reservation.');
                        }
                    });
                });
            } else {
                rentalHistoryList.innerHTML = '<p>No rental history available.</p>';
            }
        } else {
            alert('Failed to fetch rental history!');
        }
    } catch (error) {
        console.error('Error fetching rental history:', error);
        alert('An error occurred while fetching rental history.');
    }
});

    } catch (error) {
        console.error('Error:', error);
        alert('An error occurred while loading the page.');
    }
});
