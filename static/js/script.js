var app = new Vue({
    el: '#app',
    template: `
      <div>
        {{ message }}
        <form id="chatbox">
            <textarea></textarea>
            <input type="submit" value="Send">
        </form>
      </div>`,
    data: {
      message: 'Привет, Vue!'
    }
});
