let editor = document.querySelector("#editor");

window.onload = function() {
  editor = ace.edit("editor");
  editor.setTheme("ace/theme/dracula");
  editor.session.setMode("ace/mode/c_cpp");
}
function changeLang() {

    let language = document.getElementById("languages").value;

    if(language == 'c' || language == 'cpp')editor.session.setMode("ace/mode/c_cpp");
    else if(language == 'python')editor.session.setMode("ace/mode/python");
    else if(language == 'go')editor.session.setMode("ace/mode/golang");
}

