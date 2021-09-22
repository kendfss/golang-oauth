const format = function(string) {
  var args = Array.from(arguments).slice(1);
  return string.replace(/{(\d+)}/g, function(match, number) { 
    return typeof args[number] != 'undefined'
      ? args[number]
      : match;
  });
};
const findPort = function() {
  if (mode == "dev") {
    let parts = window.location.origin.split(':');
    return parts[parts.length - 1];
  }
};

const redirects = {
  "dev": "http://localhost:"+findPort()+"/oauth/redirect",
  "pro": ""
};
const mode = (1) ? "dev" : "pro";
const creds = JSON.parse(credata);
console.log(creds);
const github = document.querySelector("#github");
github.href = format(github.href, creds.github.clientID, redirects[mode]);
const google = document.querySelector("#google");
google.href = format(google.href, creds.google.clientID, redirects[mode]);
