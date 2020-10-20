$(function(){
    if (!window.EventSource) {
        alert("No EventSource!");
        return;
    }

    var $chatlog = $("#chat-log");
    var $chatmsg = $("#chat-msg");

    var isBlank = function (string) {
        return string == null || string.trim() === "";
    };

    var username;
    
    while (isBlank(username)) {
        username = prompt("What's your name");
        if (!isBlank(username)) {
            $("#chat-name").html('<b>' + username + '</b>');
        }
    }

    $("#input-form").on("submit", function (e) {
        
        $.post("/messages", {
            msg: $chatmsg.val(),
            name: username
        });
        $chatmsg.val("");
        $chatmsg.focus(); 
        // 다른 페이지로 넘어가지 않게 false return
        return false;
    });

    // 서버가 이벤트소스를 통해 메시지를 알려준다.
    var addMessage = function (data) {
        var text = "";
        if (!isBlank(data.name)) {
            text = '<strong>' + data.name + ':</strong>';
        }
        text += data.msg;
        $chatlog.prepend('<div><span>' + text + '</span></div>');
    };

    // Event Source를 연다.
    var es = new EventSource('/stream')
    es.onopen = function (e) {
        $.post('users/', {name: username})
    }

    es.onmessage = function (e) {
        var msg = JSON.parse(e.data)
        console.log(msg,"msg")
        addMessage(msg)
    }

    // 윈도우가 닫히기 직전에 발생
    window.onbeforeunload = function () {
        $.ajax({
            url: "/users?username=" + username,
            type: "DELETE"
        });
        es.close()
    };
})