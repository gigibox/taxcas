webpackJsonp([2],{"4Pe/":function(t,e){},HXef:function(t,e,s){"use strict";Object.defineProperty(e,"__esModule",{value:!0});var a={name:"Home",data:function(){return{num:0,cerList:[],curCode:"",openid:"",token:""}},created:function(){this.getOpenid()},methods:{getOpenid:function(){var t=this,e=this.GetQueryString("code");window.location.href;null==e?window.location.href="https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx6978819270b1b14b&redirect_uri=https://tax.caishuidai.com/m&response_type=code&scope=snsapi_userinfo&state=access#wechat_redirect":(t.curCode=e,t.$axios.get(t.baseUrl+"/api/weixin/openid/"+t.curCode).then(function(e){t.openid=e.data.data.openid,t.token=e.data.data.token,localStorage.setItem("openid",t.openid),t.getCerList()}))},getUserData:function(t){var e=this;e.$axios.get(e.baseUrl+"/api/v1/weixin/users/"+e.openid).then(function(s){if(s.data.success){var a=s.data.data.Certs,i=t;for(var r in a)for(var n in a[r]){var c=n,o=a[r][n];for(var u in i)i[u].cert_id==c&&(i[u].cert_status=o)}console.log(i),e.cerList=i}else e.cerList=t,console.log(s.data.msg)})},getCerList:function(){var t=this;t.$axios.get(t.baseUrl+"/api/v1/weixin/certs").then(function(e){e.data.success?t.getUserData(e.data.data):t.$toast.center("数据获取失败")})},increase:function(){this.num++,this.setNum(this.num)},clickList:function(t,e,s){0==e?this.getOrderMsg(t):5==e?this.$router.push({path:"/lookcert"}):this.$router.push({path:"/SetForm",query:{curcertid:t,curstatus:e,curcertname:s}})},GetQueryString:function(t){var e=new RegExp("(^|&)"+t+"=([^&]*)(&|$)"),s=window.location.search.substr(1).match(e);return null!=s?unescape(s[2]):null},getOrderMsg:function(t){var e=this,s=this,a=localStorage.getItem("openid");s.$axios({method:"get",url:s.baseUrl+"/api/v1/weixin/wxorder/"+a+"/"+t}).then(function(t){if(t.data.success){var a=t.data.data,i=a.appid,r=a.prepay_id,n=a.price,c=a.name,o=a.apikey,u=a.orderid;"0"===r?e.$router.push({path:"/Wait"}):e.$router.push({path:"/Order",query:{appid:i,prepayid:r,price:n,name:c,orderid:u,apikey:o}})}else s.$toast.center(t.data.msg)})}},computed:{}},i={render:function(){var t=this,e=t.$createElement,s=t._self._c||e;return s("div",{staticClass:"list"},t._l(t.cerList,function(e){return s("div",{staticClass:"list-item"},[s("div",{staticClass:"item-left"},[t._m(0,!0),t._v(" "),s("strong",{staticClass:"item-tit"},[t._v(t._s(e.cert_name))])]),t._v(" "),s("div",{staticClass:"item-right",attrs:{"data-status":e.cert_status},on:{click:function(s){t.clickList(e.cert_id,e.cert_status,e.cert_name)}}},[2==e.cert_status?s("span",{staticClass:"msg-default"},[t._v("领取中")]):3==e.cert_status?s("span",{staticClass:"msg-default msg-err"},[t._v("领取中")]):0==e.cert_status?s("span",{staticClass:"msg-default msg-err"},[t._v("未支付")]):4==e.cert_status?s("span",{staticClass:"msg-default msg-err"},[t._v("领取失败")]):5==e.cert_status?s("span",{staticClass:"msg-default msg-primay"},[t._v("查看证书")]):6==e.cert_status?s("span",{staticClass:"msg-default msg-err"},[t._v("领取失败")]):7==e.cert_status?s("span",{staticClass:"msg-default msg-err"},[t._v("领取失败")]):s("span",{staticClass:"msg-default"},[t._v("申请领取")]),t._v(" "),t._m(1,!0)])])}))},staticRenderFns:[function(){var t=this.$createElement,e=this._self._c||t;return e("span",{staticClass:"img-cont"},[e("img",{attrs:{src:s("kkb8"),alt:""}})])},function(){var t=this.$createElement,e=this._self._c||t;return e("span",{staticClass:"img-icon"},[e("img",{attrs:{src:s("qNF3"),alt:""}})])}]};var r=s("C7Lr")(a,i,!1,function(t){s("4Pe/")},"data-v-7eaaf7ac",null);e.default=r.exports},kkb8:function(t,e,s){t.exports=s.p+"static/img/index-icon.43c52de.png"},qNF3:function(t,e){t.exports="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAwAAAAVCAYAAAEFqz8iAAAABGdBTUEAALGPC/xhBQAAAThJREFUOBGlUkFugzAQtI0/kHuVBrUvyKGHnoKAx/bQFwCCSw455AWp2iR9Qq+AO+uwjnEiEilIZmdnZ2dt2aIoik/BHxLD2I+Od4CrlsDvgwmBZM6JIsCE9BNNida6pyDqun7CWhBJyQkAy6H5APyc5/kR0X7WjFAocIWTToiqqt77vl/biUTSsLZtf5VScZqmUjOBWowZZ4eyLLds40dljJlh8LdPOkyFx4vnXQzGsDRSysbNIQByh13ub5NQfV0or5LwogP+ZVm2GPkikejYorgE/sFa+XcSiim3x8Au58ANVow12Tg69z2Nowa422+q8WrDVKN9xSwI4/CGO+JxPQq5ck/NFw+vrMbTe4XwEEVRnCTJnjSjBhJ2XddA+BIK2dA2hEKMdo4s5KhxcRs4vpHjlJAb/gH8iMrCUgJxyQAAAABJRU5ErkJggg=="}});
//# sourceMappingURL=2.e16c277a083446980155.js.map