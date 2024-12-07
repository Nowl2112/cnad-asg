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
                }
            } catch (error) {
                console.error('Error updating profile:', error);
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
                        
                                // Fetch reservation details from the backend
                                try {
                                    const response = await fetch(`http://localhost:8082/reservations/${reservationId}`);
                                    if (response.ok) {
                                        const reservation = await response.json();
                        
                                        // Pre-fill the form with reservation data
                                        document.getElementById('editStartTime').value = reservation.startTime
                                        document.getElementById('editEndTime').value = reservation.endTime
                        
                                        // Store reservation ID in a hidden attribute for the form
                                        document.getElementById('editReservationForm').dataset.reservationId = reservationId;
                        
                                        // Show the popup
                                        document.getElementById('editReservationPopup').style.display = "flex";
                                    } else {
                                        alert('Failed to fetch reservation details!');
                                    }
                                } catch (error) {
                                    console.error('Error fetching reservation details:', error);
                                    alert('An error occurred while fetching reservation details.');
                                }
                            });
                        });
                        document.getElementById('cancelReservation').addEventListener('click', async () => {
                            const reservationId = document.getElementById('editReservationForm').dataset.reservationId;
                        
                            if (confirm('Are you sure you want to cancel this reservation?')) {
                                try {
                                    const response = await fetch(`http://localhost:8082/reservations/${reservationId}/cancel`, {
                                        method: 'PUT',
                                    });
                        
                                    if (response.ok) {
                                        alert('Reservation canceled successfully!');
                                        document.getElementById('editReservationPopup').style.display = "none";
                                        document.getElementById('viewHistory').click(); // Refresh the rental history
                                    } else {
                                        alert('Failed to cancel reservation!');
                                    }
                                } catch (error) {
                                    console.error('Error canceling reservation:', error);
                                    alert('An error occurred while canceling the reservation.');
                                }
                            }
                        });
                     
                        document.getElementById('editReservationForm').addEventListener('submit', async (e) => {
                            e.preventDefault();
                        
                            function formatLocalTime(dateTimeLocal) {
                                const localTime = new Date(dateTimeLocal);
                                
                                // Manually construct the ISO format string based on local time
                                const year = localTime.getFullYear();
                                const month = String(localTime.getMonth() + 1).padStart(2, '0'); // Months are 0-indexed
                                const day = String(localTime.getDate()).padStart(2, '0');
                                const hours = String(localTime.getHours()).padStart(2, '0');
                                const minutes = String(localTime.getMinutes()).padStart(2, '0');
                                const seconds = String(localTime.getSeconds()).padStart(2, '0');
                                
                                return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}Z`;
                            }
                        
                            const reservationId = e.target.dataset.reservationId;
                            const newStartTime = document.getElementById('editStartTime').value;
                            const newEndTime = document.getElementById('editEndTime').value;
                            
                            const startTimeFormatted = formatLocalTime(newStartTime);
                            const endTimeFormatted = formatLocalTime(newEndTime);
                        
                            // Ensure the new times are valid
                            if (!newStartTime || !newEndTime) {
                                alert('Please provide valid start and end times.');
                                return;
                            }
                        
                            try {
                                // Fetch the current reservation status before updating
                                const statusResponse = await fetch(`http://localhost:8082/reservations/${reservationId}`);
                                if (!statusResponse.ok) {
                                    alert('Failed to retrieve current reservation details.');
                                    return;
                                }
                                const currentReservation = await statusResponse.json();
                        
                                // Update reservation with the correct payload format
                                const response = await fetch(`http://localhost:8082/reservations/${reservationId}`, {
                                    method: 'PUT',
                                    headers: { 'Content-Type': 'application/json' },
                                    body: JSON.stringify({
                                        start_time: startTimeFormatted,  // Change the key to start_time
                                        end_time: endTimeFormatted,      // Change the key to end_time
                                        status: currentReservation.status,
                                    }),
                                });
                        
                                if (response.ok) {
                                    alert('Reservation updated successfully!');
                                    document.getElementById('editReservationPopup').style.display = "none";
                                    document.getElementById('viewHistory').click(); // Refresh the rental history
                                } else {
                                    alert('Failed to update reservation!');
                                }
                            } catch (error) {
                                console.error('Error updating reservation:', error);
                                alert('An error occurred while updating the reservation.');
                            }
                        });
                        
                        
                        document.querySelectorAll('.complete-btn').forEach(button => {
                            button.addEventListener('click', async (e) => {
                                const reservationId = e.target.dataset.id;
                                
                                // Find the corresponding item to get the total price
                                const item = Array.from(rentalHistoryList.children).find(listItem =>
                                    listItem.querySelector(`.complete-btn[data-id="${reservationId}"]`)
                                );
                                
                                if (item) {
                                    let totalPrice = parseFloat(item.querySelector('p:nth-child(4)').textContent.split('$')[1]); // Convert to dollars
                                    
                                    // Ensure totalPrice is an integer (in cents)
                                    totalPrice = Math.round(totalPrice * 100);  // Convert to cents and round to the nearest integer
                                    
                                    // Create an array of items
                                    const items = [{
                                        id: reservationId,  // Passing reservation ID as Item ID
                                        amount: totalPrice   // Amount is now guaranteed to be an integer
                                    }];
                                    
                                    try {
                                        // Call backend to create a Payment Intent
                                        const response = await fetch('http://localhost:8083/create-payment-intent', {
                                            method: 'POST',
                                            headers: { 'Content-Type': 'application/json' },
                                            body: JSON.stringify({ items }), // Pass the array of items
                                        });
                                
                                        if (response.ok) {
                                            const { clientSecret } = await response.json();
                                
                                            // Redirect to the payment page with the client secret and reservation ID
                                            const url = `payment.html?client_secret=${clientSecret}&reservation_id=${reservationId}`;
                                            window.location.href = url;
                                        } else {
                                            alert('Failed to generate payment intent.');
                                        }
                                    } catch (error) {
                                        console.error('Error creating payment intent:', error);
                                        alert('An error occurred while initiating payment.');
                                    }
                                } else {
                                    console.error('Item data not found for reservation ID:', reservationId);
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
