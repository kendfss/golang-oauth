const format = function(string) {
  var args = Array.from(arguments).slice(1);
  return string.replace(/{(\d+)}/g, function(match, number) { 
    return typeof args[number] != 'undefined'
      ? args[number]
      : match
    ;
  });
};

const creds = JSON.parse(credata);
const redirects = {
  "dev": "http://localhost:8080/oauth/redirect",
  "pro": ""
};
const mode = (1) ? "dev" : "pro";
const github = document.querySelector("#github");
github.href = format(github.href, creds.gh.clientID, redirects[mode]);
const google = document.querySelector("#google");
google.href = format(google.href, creds.goo.clientID, redirects[mode]);
