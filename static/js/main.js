$(document).ready(function() {
  GetVpcs();
});

var GetVpcs = function() {
  const data = {action: "getvpcs"};
  request(data, (res)=>{
    console.log(res);
    if (!!res && !!res.message && res.message.length > 0) {
      vpcData = JSON.parse(res.message);
    }
    $("#result1").text(JSON.stringify(vpcData, null, "\t"));
    $("#info").removeClass("hidden").addClass("visible");
    GetSubnets();
  }, (e)=>{
    console.log(e.responseJSON.message);
    $("#warning").text(e.responseJSON.message).removeClass("hidden").addClass("visible");
  });
};

var GetSubnets = function() {
  const data = {action: "getsubnets"};
  request(data, (res)=>{
    console.log(res);
    if (!!res && !!res.message && res.message.length > 0) {
      subnetData = JSON.parse(res.message);
    }
    $("#result2").text(JSON.stringify(subnetData, null, "\t"));
  }, (e)=>{
    console.log(e.responseJSON.message);
    $("#warning").text(e.responseJSON.message).removeClass("hidden").addClass("visible");
  });
};

var request = function(data, callback, onerror) {
  $.ajax({
    type:          'POST',
    dataType:      'json',
    contentType:   'application/json',
    scriptCharset: 'utf-8',
    data:          JSON.stringify(data),
    url:           App.url
  })
  .done(function(res) {
    callback(res);
  })
  .fail(function(e) {
    onerror(e);
  });
};
var App = { url: location.origin + {{ .ApiPath }} };
