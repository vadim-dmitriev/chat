# Пользовательские сценарии

## Отправка сообщения существующему пользователю

Некий человек узнал о существовании мессенджера от своего знакомого. Этот знакомый уже является пользователем чата и зарегистрирован в нем. Человек хочет так же стать пользователем мессенджера для того, чтобы общаться со свои знакомым.

### Этапы
 
 - Регистрация в мессенджере

Для регистрации необходимо придумать имя пользователя и пароль. После этого ввести эти данные в соответствующие поля и нажать `Sign up`. Имя пользователя должно быть уникально.

 - Аутентификация в мессенджере

Для аутентификации необходимо ввести свои имя пользователя и пароль после чего нажать кнопку `Sign in`. Если эти данные совпадут, то пользователь попадет на страницу чата, в котором еще нет никаких диалогов и чатов. В противном случае пользователю будет сказано, что он ввел неверные данные.

 - Поиск пользователя для отправки сообщения

Для того, чтобы найти пользователя внутри мессенджера необходимо знать его имя. В поле поиска пользователей, после нажатия соответствующей кнопки, нужно ввести необходимое имя. Если пользователь с таким именем зарегистрирован, то создается пустой (без сообщений) диалог. В случае, если такого пользователя не нашлось или искомый пользователь это он сам (ищу себя), диалог создан не будет.

 - Отправка сообщения 

Для того, чтобы отправить сообщение, необходимо сперва выбрать диалог. В нашем случае он был создан на предыдущем этапе. После чего, набрать необходимое сообщение и нажать кнопку отрпавить. Сообщение будет отправленно в любом случае, вне зависимости от того в сети сейчас собеседник или нет.