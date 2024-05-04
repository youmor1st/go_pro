$(document).ready(function() {
    $('#register-form').submit(function(event) {
        event.preventDefault();

        var username = $('#username').val();
        var password = $('#password').val();
        var email = $('#email').val();

        $.ajax({
            type: 'POST',
            url: '/register',
            contentType: 'application/json',
            data: JSON.stringify({
                username: username,
                password: password,
                email: email
            }),
            success: function(response) {
                $('#message').text('Registration successful!');
            },
            error: function(xhr, status, error) {
                $('#message').text(xhr.responseText);
            }
        });
    });
});
