Vue.component("conversation", {
	data: function () {
		return {}
	},
	props: ["index", "conversation", "isChoosed"],
	template: `
		<div class="conversation" v-bind:class="{ conversationClicked: this.isChoosed }">
			<strong>{{ conversation.name }}</strong>
			<br/><br/>
			{{ conversation.lastMessage }}
		</div>`,
	methods: {
		select: function() {
			this.isChoosed = true
			this.conversation.lastMessage = "asd"
		},
	}
})

Vue.component("search-user-panel", {
	data: function () {
		return {
			username: ""
		}
	},
	props: ["isActive"],
	template: `
		<div v-show="isActive">
			<button v-on:click="hidePanel">x</button>
			<input v-model="username" />
			<button v-on:click="searchUser">Search</button>
		</div>
	`,
	methods: {
		hidePanel: function() {
			this.$emit("hidePanel");
		},
		searchUser: function() {
			if (this.username === "") {
				return
			}

			this.$emit("search-user", this.username)
			this.username = ""
			this.hidePanel()
		}
	}
})

Vue.component("user-panel", {
	data: function () {
		return {
			isUserSearchPanelActive: false,
		}
	},
	template: `
		<div class="user-panel">
			<div style="border: 1px black solid;" v-on:click="showSearchUserPanel">
				<p>Найти пользователя</p>
			</div>
			<search-user-panel v-bind:isActive="isUserSearchPanelActive"
							   v-on:hidePanel="hideSearchUserPanel"
							   v-on:search-user="searchUser"/>
		</div>`, 
	methods: {
		showSearchUserPanel: function() {
			this.isUserSearchPanelActive = true;
		},
		hideSearchUserPanel: function() {
			this.isUserSearchPanelActive = false;
		},
		searchUser: function(username) {
			this.$emit("search-user", username)
		}
	}
})

Vue.component("conversations", {
	data: function () {
		return {
			isChoosed: new Array(this.conversations.length).fill(false),
			choosedConversation: -1
		}
	},
	props: ["conversations"],
	template: `
		<div class="conversations">
			<div v-for="(conversation, index) in conversations">
				<conversation v-bind:isChoosed="isChoosed[index]" v-bind:conversation="conversation" v-bind:index="index" v-on:click.native="change(index)"/>
			</div>
		</div>`,
	methods: {
		change: function(index) {
			if (this.choosedConversation != index) {
				this.isChoosed = new Array(this.conversations.length).fill(false);
				this.isChoosed[index] = true;
			}
		}
	},
	watch: {
		conversations: function() {
			this.change(this.conversations.length-1)
		}
	}
})

Vue.component("chat", {
	data: function() {
		return {}
	},
	props: ["messages"],
	template: `
		<div class="chat">
			<ul>
				<li v-for="message in messages">
					{{ message }}
				</li>
			</ul>
			<div id="control">
				<input v-model="currentMessage" v-on:keyup.enter="sendMessage"/>
				<button v-on:click="sendMessage" value="Send">Send</button>
			</div>
		</div>`,
	methods: {
		sendMessage: function() {
			// if (this.currentMessage !== "") {
			// 	this.messages.push(this.currentMessage);
			// 	this.currentMessage = "";
			// }
		},
	}
})

var app = new Vue({
	el: '#app',
	data: {
		conversations: [
			{name: "Витя", lastMessage: "Добрый день, ок"},
			{name: "Сережа", lastMessage: "Привет! Во сколько ты приедешь?"},
			{name: "Веррроника", lastMessage: "kek"},
			{name: "Павлик", lastMessage: "Отвечаю на вопросы которые задает сегодня..."},
		],
		currentMessage: "",
		messages: [],
		ws: null,
		},

	created: function() {
		const conversations = this.conversations

		this.ws = new WebSocket("ws://"+window.location.host+"/api/v1/ws");

		this.ws.onmessage = function(event) {
			message = JSON.parse(event.data)

			switch (message.action) {
			case "newConversation":
				conversations.push({
					"name": message.name,
					"lastMessage": message.lastMessage,
				});
				break;
			}
		}
	},
	destroyed: function() {
		this.ws.close()
	},
	methods: {
		searchUser: function(username) {
			// alert(username);
			this.ws.send(
				JSON.stringify({
					"action": "searchUser",
					"username": username
				})
			)
		}
	}
});
