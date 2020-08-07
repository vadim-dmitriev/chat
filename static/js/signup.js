var app = new Vue({
    el: "#app",
    data: {
        login: "",
        password: "",
    },
    methods: {
        doRegister: function() {
            fetch("/api/v1/register", {
                method: "POST",
                headers: {
                    'Content-Type': 'application/json;charset=utf-8'
                },
                body: JSON.stringify({
                    "login": this.login,
                    "password": this.password,
                }),
            }).then(
                response => response.text().then(
                    text => {
                        if (text === "ok") {
                            window.location.replace("/signin");
                        }
                    }
                ))
        }
    }

})