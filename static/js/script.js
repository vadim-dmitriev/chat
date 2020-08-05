Vue.component("conversation", {
	data: function () {
		return {}
	},
	props: ["index", "conversation", "choosedConversation"],
	template: `
		<div class="conversation" v-on:click="select" v-bind:class="{ conversationClicked: conversation.isChoosed }">
			{{ conversation.lastMessage }}
			{{ conversation.isChoosed }}
		</div>`,
	methods: {
		select: function() {
			this.$emit("select", this.index)
			if (this.choosedConversation == this.index) {
				this.conversation.isChoosed = true
				console.log("here")
			} else {
				this.conversation.isChoosed = false
			}
		},
	}
})

Vue.component("user-panel", {
	data: function () {
		return {}
	},
	template: `
		<div class="user-panel">
		</div>`, 
})

Vue.component("conversations", {
	data: function () {
		return {
			choosedConversation: -1
		}
	},
	props: ["conversations"],
	template: `
		<div class="conversations">
			<div v-for="(conversation, index) in conversations">
				<conversation v-bind:choosedConversation="choosedConversation" v-bind:conversation="conversation" v-bind:index="index" v-on:select="select"/>
			</div>
		</div>`,
	methods: {
		select: function(newChoosedConversationIndex) {
			console.log(newChoosedConversationIndex)
			this.choosedConversation = newChoosedConversationIndex
		}
	}
})

new Vue({
	el: '#app',
	template: `
	<div>
		<user-panel />
		<conversations v-bind:conversations="conversations"/>
		<div id="chat">
		<ul>
			<li v-for="message in messages">
			{{ message }}
			</li>
		</ul>
		<div id="control">
			<input v-model="currentMessage"
				v-on:keyup.enter="sendMessage"/>
			<button v-on:click="sendMessage" value="Send">Send</button>
		</div>
		</div>
	</div>`,
	data: {
	conversations: [
		{lastMessage: "Добрый день, ок", isChoosed: false},
		{lastMessage: "Привет! Во сколько ты приедешь?", isChoosed: false},
		{lastMessage: "Отвечаю на вопросы которые задает сегодня...", isChoosed: false},
		{lastMessage: "kek", isChoosed: false},
		{lastMessage: "Поздравляем, Ваш номер подтвержден!", isChoosed: false},
		{lastMessage: "Any other questions? Something didn`t...", isChoosed: false},
		{lastMessage: "Mac Miller - Dunno", isChoosed: false},
		{lastMessage: "Отвечаю на вопросы которые задает сегодня...", isChoosed: false},
	],
	currentMessage: "",
	messages: [],
	connection: null,
	},
	created: function() {
	this.connection = new WebSocket("ws://"+window.location.host+"/ws");
	},
	methods: {
	sendMessage: function() {
		if (this.currentMessage !== "") {
		this.messages.push(this.currentMessage);
		this.currentMessage = "";
		}
	},
	}

});
