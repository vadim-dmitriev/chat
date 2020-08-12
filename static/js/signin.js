var app = new Vue({
    el: "#app",
    data: {
        login: "",
        password: "",
    },
    methods: {
        doAuth: function() {
            var xhr = new XMLHttpRequest();

            xhr.open("POST", "/api/v1/auth", false)

            // xhr.send()
            xhr.send(JSON.stringify({
                "login": this.login,
                "password": this.password,
            }));

            if (xhr.status == 200) {
                window.location.replace("/");
            }

        
        }
    }
})