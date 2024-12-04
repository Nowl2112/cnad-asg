document.addEventListener('DOMContentLoaded', async function () {
    const userId = localStorage.getItem('user_id');
    if (!userId) {
        alert("Not logged in")
        window.location.href = 'index.html'; // Redirect to login if not logged in
        return;
    }

    const response = await fetch(`http://localhost:8080/user/${userId}`);
    if (response.ok) {
        const user = await response.json();
        document.getElementById('email').value = user.email;
        document.getElementById('phone').value = user.phone;
        document.getElementById('membership_tier').value = user.membership_tier;

    } else {
        alert('Failed to fetch user data!');
    }
});

document.getElementById('profileForm').addEventListener('submit', async function (e) {
    e.preventDefault();
    const userId = localStorage.getItem('userId');
    const phone = document.getElementById('phone').value;

    const response = await fetch(`http://localhost:8080/user/${userId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ phone }),
    });

    if (response.ok) {
        alert('Profile updated successfully!');
    } else {
        alert('Failed to update profile!');
    }
});
