var app = new Vue({
    el: "#app",
    data: {
        login: "",
        password: "",
        helpMessage: "",
    },
    methods: {
        doAuth: function() {
            if (this.login == '') {
                this.showHelpMessage('Введите ');
            }
            fetch('/api/v1/signin', {
                method: 'POST',
                body: JSON.stringify({
                    "username": this.login,
                    "password": this.password,
                })
            }).then(resp => {
                if (resp.status == 200) {
                    window.location.replace('/');
                } else {
                    if (resp.status == 401) {
                        this.showHelpMessage('Неверный логин или пароль');
                    }
                }
            })
        },
        showHelpMessage: function(message) {
            this.helpMessage = message;
            this.isShowHelpMessage = true;

        }
    },
    computed: {
        signInDisabled: function() {
            if (this.login.length==0 || this.password.length==0) {
                return true
            }
            return false
        },
        signInClass: function() {
            if (!this.signInDisabled) {
                return "btn btn-outline-success btn-lg"
            }
            return "btn btn-outline-secondary btn-lg"
        }
    }
});