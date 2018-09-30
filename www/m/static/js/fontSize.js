// 获取屏幕宽度
let htmlWidth =
  document.documentElement.clientWidth || document.body.clientWidth;
// 获取html的dom元素
let htmlDom = document.getElementsByTagName("html")[0];
if (htmlWidth > 750) {
  htmlWidth = 750;
}
htmlDom.style.fontSize = htmlWidth / 10 + "px";
// 通过监听事件，时刻的动态的修改font-size的大小(屏幕宽度变化时,font-size的大小也随着变化)
window.addEventListener("resize", function(e) {
  // 获取屏幕宽度
  htmlWidth = document.documentElement.clientWidth || document.body.clientWidth;
  // 获取html的dom元素
  if (htmlWidth > 750) {
    htmlWidth = 750;
  }
  htmlDom.style.fontSize = htmlWidth / 10 + "px"; //设计稿的大小是基于750的宽度设计的
});
