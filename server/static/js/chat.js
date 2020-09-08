Vue.component("conversation", {
	data: function () {
		return {}
	},
	props: ["conversation", "choosedConversation"],
	template: `
		<div class="conversation" v-bind:class="{ conversationClicked: isChoosed }">
			<strong>{{ conversation.name }}</strong>
			<br/><br/>
			{{ lastMessage }}
			<br/><br/>
			{{ lastMessageTime }}
		</div>`,
	methods: {},
	computed: {
		isChoosed: function() {
			return this.conversation.name === this.choosedConversation
		},
		lastMessage: function() {
			if (this.conversation.messages == null || this.conversation.messages.length == 0) {
				return "Нет сообщений"
			}
			return this.conversation.messages[0].text
		},
		lastMessageTime: function() {
			if (this.conversation.messages == null || this.conversation.messages.length == 0)  {
				return ""
			}
			return this.conversation.messages[0].datetime
		}

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
			choosedConversation: ""
		}
	},
	props: ["conversations"],
	template: `
		<div class="conversations">
			<div v-for="conversation of conversations">
				<conversation :conversation="conversation"  v-on:click.native="change(conversation.name)"
							  :choosedConversation="choosedConversation"
				/>
			</div>
		</div>`,
	methods: {
		change: function(convName) {
			if (this.choosedConversation !== convName) {
				this.choosedConversation = convName;
				this.$emit("change-conversation", convName);
			}
		}
	},
});

Vue.component("message", {
	data: function() {
		return {}
	},
	props: ["message"],
	template: `
		<div class="message" :align="msgType">
		{{ message.text }}
			
		</div>
	`,
	computed: {
		msgType: function () {
			if (this.message.from.id == USER) {
				return "right"
			}
			return "left"
		},
		rectHeight: function() {
			if (this.message.text.length < 10) {
				return 30
			}
			if (String(this.message.text).includes(" ") && this.message.text.length > 30)   {
				return this.message.text.length 
			}
			return 30
		},
		rectWidth: function() {
			if (this.message.text.length  > 10) {
				return 300
			}
			return this.message.text.length  * 30
		}
	}
})

Vue.component("chat", {
	data: function() {
		return {
			currentMessage: "",
		}
	},
	props: ["conversation"],
	template: `
		<div class="chat" v-show="isActive">
			<ul class="messages">
				<li v-for="message in reversedMessages">
					<message :message="message" />
				</li>
			</ul>
			<div>
				<input v-model="currentMessage" v-on:keyup.enter="sendMessage"/>
				<button v-on:click="sendMessage" value="Send">Send</button>
			</div>
		</div>
	`,
	methods: {
		sendMessage: function() {
			if (this.currentMessage !== "") {
				this.$emit("send-message", this.currentMessage)
				this.conversation.messages.unshift({
					value: this.currentMessage,
					time: new Date().toString(),
				})
				this.currentMessage = ""
			}
		},
	},
	computed: {
		reversedMessages: function() {
			if (this.conversation == null || this.conversation.messages == null || this.conversation.messages.length == 0) {
				return []
			}
			return this.conversation.messages.slice().reverse()
		},
		messages: function() {
			if (this.conversation == null || this.conversation.messages == null || this.conversation.messages.length == 0) {
				return {
					value: "Нет сообщений"
				}
			}
			return this.conversation.messages
		},
		isActive: function() {
			return (this.conversation != null)
		}
	},
});

let USER = null;

var app = new Vue({
	el: '#app',
	data: {
		conversations: [],
		currentConversation: null,
		currentConversationName: "",
		ws: null,
	},

	created: function() {
		t = this

		this.ws = new WebSocket("ws://"+window.location.host+"/api/v1/ws");

		this.ws.onopen = function() {
			t.getUser();
			t.getConversaions();
		}

		this.ws.onmessage = function(event) {
			message = JSON.parse(event.data)
			console.log("NEW MESSAGE FROM SERVER:", message);

			switch (message.action) {
			case "getUser":
				USER = message.data.userID;
				break;

			case "getConversations":
				t.conversations = message.data;
				break;

			case "searchUser":
				if (message.isUserExists) {
					t.conversations.push({
						name: message.newConversationWith,
						is_dialog: true,
						messages: [
							{
								value: "Нет сообщений...",
								sender: "",
								time: "",
							}
						]
					})
					t.currentConversation = t.conversations[t.conversations.length -1]
					t.currentConversationName = message.newConversationWith
				}
				break;

			case "newMessage":
				for (conversation of t.conversations) {
					if (conversation.id === message.data.to.id) {
						conversation.messages.unshift({
							text: message.data.text,
							from: {
								id: message.data.from.id
							},
							to: {
								id: message.data.to.id,
							},
							time: new Date().toString()
						});
						return;
					}
				}
				t.conversations.push({
					name: message.data.to.name,
					is_dialog: true,
					messages: [
						{
							text: message.data.text,
							from: {
								id: message.data.from.id
							},
							to: {
								id: message.data.to.id,
							},
							time: new Date().toString()
						},
					]
				})
				break;
			}
		}

	},
	beforeDestroy: function() {
		alert("close")
		this.ws.close()
	},
	methods: {
		getUser: function() {
			this.ws.send(
				JSON.stringify({
					action: "getUser"
				})
			);
		},
		getConversaions: function () {
			this.ws.send(
				JSON.stringify({
					action: "getConversations"
				})
			);

		},
		searchUser: function(username) {
			for (conversation of this.conversations) {
				if (conversation.name == username) {
					this.currentConversation = conversation
					this.currentConversationName = conversation.name
					return
				}
			}
			this.ws.send(
				JSON.stringify({
					action: "searchUser",
					username: username
				})
			);
		},
		sendMessage: function(message) {
			this.ws.send(
				JSON.stringify({
					action: "newMessage",
					data: {
						conversationID: this.currentConversation.id,
						text: message
					}
				})
			);
		},
		changeConversation: function(currentChoosedConversation) {
			this.currentConversationName = currentChoosedConversation

			for (conversation of this.conversations) {
				if (conversation.name == currentChoosedConversation) {
					this.currentConversation = conversation
				}
			}
			
		}
	}
});
