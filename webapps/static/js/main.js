const API_BASE_URL = 'http://localhost:8081';
document.addEventListener('DOMContentLoaded', function () {
    const userForm = document.getElementById('userForm');
    if (userForm) {
        userForm.addEventListener('submit', async function (event) {
            event.preventDefault();

            const submitButton = event.submitter;
            const method = submitButton.dataset.method;
            const action = submitButton.dataset.action;

            const formData = new FormData(userForm);
            const data = {};
            formData.forEach((value, key) => {
                data[key] = value;
            });

            const payload = {
                name: data.name,
                email: data.email,
                phone_number: data.phone_number,
                address: {
                    country: data.country,
                    state: data.state,
                },
            };

            if (data.dob) {
                payload.dob = new Date(data.dob).toISOString();
            }

            try {
                const response = await fetch(`${API_BASE_URL}${action}`, {
                    method: method,
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(payload),
                });

                if (!response.ok) {
                    const errorData = await response.json();
                    throw new Error(errorData.message || 'Something went wrong');
                }

                alert('Operation successful!');
                window.location.href = '/users';
            } catch (error) {
                console.error('Error submitting form:', error);
                alert(`Error: ${error.message}`);
            }
        });
    }
});

async function deleteUser(userId) {
    if (!confirm('Are you sure you want to delete this user?')) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/v1/users/${userId}`, {
            method: 'DELETE',
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.message || 'Failed to delete user');
        }

        alert('User deleted successfully!');
        window.location.reload();
    } catch (error) {
        console.error('Error deleting user:', error);
        alert(`Error: ${error.message}`);
    }
}
