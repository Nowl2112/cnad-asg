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
            document.getElementById('editProfilePopup').style.display="flex";
        });

        // Event listener for closing the edit profile popup
        document.querySelector('.close-btn').addEventListener('click', () => {
            document.getElementById('editProfilePopup').style.display="none";
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
                    document.getElementById('editProfilePopup').classList.add('hidden');
                } else {
                    alert('Failed to update profile!');
                }
            } catch (error) {
                console.error('Error updating profile:', error);
                alert('An error occurred while updating your profile.');
            }
        });

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
                                <hr>
                            `;
                            rentalHistoryList.appendChild(listItem);
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
