webpackJsonp([3],{"2/T4":function(e,a,t){"use strict";var s={render:function(){var e=this,a=e.$createElement,t=e._self._c||a;return t("div",{},[t("div",{staticClass:"handle-box"},[t("el-form",{staticClass:"inlineform1",attrs:{inline:!0}},[t("el-form-item",{attrs:{label:""}},[t("el-button",{staticClass:"handle-del mr10",attrs:{type:"primary",icon:"plus"},on:{click:e.out}},[e._v("批量导出")])],1)],1),e._v(" "),t("el-form",{staticClass:"inlineform2",attrs:{inline:!0}},[t("el-form-item",{attrs:{label:""}},[t("el-input",{staticClass:"handle-input mr10",attrs:{placeholder:"请输入渠道名称或账号查找"},model:{value:e.searchWord,callback:function(a){e.searchWord=a},expression:"searchWord"}})],1),e._v(" "),t("el-form-item",[t("el-button",{attrs:{type:"primary",icon:"search"},on:{click:e.search}},[e._v("查询")])],1)],1)],1),e._v(" "),t("el-table",{directives:[{name:"loading",rawName:"v-loading",value:e.loading,expression:"loading"}],ref:"multipleTable",staticStyle:{width:"100%"},attrs:{data:e.waitlistData,border:""}},[t("el-table-column",{attrs:{type:"index",label:"编号",width:"50"}}),e._v(" "),t("el-table-column",{attrs:{prop:"cert_name",label:"申请证书",width:"180"}}),e._v(" "),t("el-table-column",{attrs:{prop:"wechatID",label:"微信ID",width:"180"}}),e._v(" "),t("el-table-column",{attrs:{prop:"name",label:"申请人"}}),e._v(" "),t("el-table-column",{attrs:{prop:"cardID",label:"身份证号",width:"180"}}),e._v(" "),t("el-table-column",{attrs:{prop:"time",label:"申请时间"}}),e._v(" "),t("el-table-column",{attrs:{prop:"money",label:"支付金额"}}),e._v(" "),t("el-table-column",{attrs:{prop:"status",label:"申请状态"}})],1),e._v(" "),t("div",{staticClass:"pagination"},[t("el-pagination",{attrs:{"current-page":e.currentPage,layout:"prev, pager, next","page-size":10,total:e.pageTotal},on:{"current-change":e.handleCurrentChange}})],1)],1)},staticRenderFns:[]};var r=t("VU/8")({data:function(){return{radio:"all",loading:!1,searchWord:"",currentPage:1,pageTotal:0,waitlistData:[]}},created:function(){},methods:{getData:function(){var e=this;e.$axios.get(e.baseUrl+"/api/v1/admin/certs").then(function(a){e.waitlistData=a.data.data})},changeTab:function(){},out:function(){},search:function(){},handleCurrentChange:function(){}}},s,!1,function(e){t("Eiw/")},"data-v-1df726f8",null);a.a=r.exports},"7a/k":function(e,a){},Dq1k:function(e,a,t){"use strict";Object.defineProperty(a,"__esModule",{value:!0});var s=t("2/T4"),r=(t("IcnI"),{components:{wExport:s.a},data:function(){return{radio:"export",loading:!1,searchWord:"",dialogImgVisible:!1,passedImg:"https://tax.caishuidai.com/api/export/image/cert1/20180901000004.png",passedCertName:"",currentPage:1,curPageSize:10,pageTotal:0,listData:[],allListData:[],curCertId:this.$route.query.cert_id,uploadUrl:this.baseUrl+"/api/v1/admin/excels",fileList:[],fileList2:[],personalIdList:[],msg:{certname:"无",name:"无",cardid:"无",date:"无",price:"无",status:"无"},showTable:!0,excelsdata:{Authorization:"Bearer "+localStorage.token},rejectDisabled:!1}},computed:{},watch:{$route:function(){this.curCertId=this.$route.query.cert_id,this.showTable=!0,this.radio="export",this.getData("export",this.currentPage,this.curPageSize)}},created:function(){this.curCertId=this.$route.query.cert_id,this.getData("export",this.currentPage,this.curPageSize)},methods:{getParams:function(){},getData:function(){var e=arguments.length>0&&void 0!==arguments[0]?arguments[0]:"export",a=arguments.length>1&&void 0!==arguments[1]?arguments[1]:1,t=arguments.length>2&&void 0!==arguments[2]?arguments[2]:10,s=this;s.$axios.get(s.baseUrl+"/api/v1/admin/applicants/certs/"+s.curCertId+"?page="+a+"&limit="+t+"&type="+e).then(function(a){a.data.success?(s.allListData=a.data.data.list,s.pageTotal=a.data.data.count,"reject"==e&&(a.data.data.list.length>0&&"0"==a.data.data.list[0].pay_amount||0==a.data.data.list.length?s.rejectDisabled=!0:s.rejectDisabled=!1)):s.$message.error(a.data.msg)})},clickDetail:function(e,a,t,s,r,i,l,n){this.showTable=!1,this.msg.certname=e,this.msg.name=a,this.msg.cardid=t,this.msg.date=s,this.msg.price=r+" 元",this.msg.status=4==l&&1==n?"未审核":i},clickDetailImg:function(e,a,t){var s=this;s.passedCertName=a,s.$axios.get(s.baseUrl+"/api/v1/admin/images/certs/"+e+"/"+t).then(function(e){e.data.success?(s.passedImg=e.data.data.image_url,s.dialogImgVisible=!0):s.$message.error(e.data.msg)})},handleSelectionChange:function(e){var a=[];for(var t in e)a.push(e[t].personal_id);this.personalIdList=a},handleSizeChange:function(e){this.curPageSize=e,this.getData(this.radio,this.currentPage,this.curPageSize)},handleCurrentChange:function(e){this.currentPage=e,this.getData(this.radio,this.currentPage,this.curPageSize)},changeTab:function(){this.showTable=!0,this.getData(this.radio)},outsExport:function(){var e=this;e.$confirm("是否批量导出全部“待导出”数据?",{confirmButtonText:"确定",cancelButtonText:"取消",type:"warning"}).then(function(){e.$axios.get(e.baseUrl+"/api/v1/admin/files/applicants/certs/"+e.curCertId,{params:{type:"export"}}).then(function(a){a.data.success?(window.location.href=a.data.data.file_url,e.$message({message:"导出成功",type:"success"}),e.getData(e.radio,e.currentPage,e.curPageSize)):e.$message.error(a.data.msg)})}).catch(function(){})},sureReject:function(){var e=this;if(e.personalIdList.length<1)return e.$message({message:"请选择数据",type:"warning"}),!1;e.$axios.put(e.baseUrl+"/api/v1/admin/applicants/certs/"+e.curCertId,{action:"reject",pids:e.personalIdList}).then(function(a){a.data.success?(e.$message({message:"操作成功",type:"success"}),e.getData(e.radio,e.currentPage,e.curPageSize)):e.$message.error(a.data.msg)})},outsReject:function(){var e=this;e.$confirm("是否批量导出全部“已拒绝”数据?",{confirmButtonText:"确定",cancelButtonText:"取消",type:"warning"}).then(function(){e.$axios.get(e.baseUrl+"/api/v1/admin/files/applicants/certs/"+e.curCertId,{params:{type:"reject"}}).then(function(a){a.data.success?(window.location.href=a.data.data.file_url,e.$message({message:"导出成功",type:"success"}),e.getData(e.radio,e.currentPage,e.curPageSize)):e.$message.error(a.data.msg)})}).catch(function(){})},insReject:function(){},search:function(){var e=this;if(""==e.searchWord)return e.$message({message:"请输入查询内容",type:"warning"}),!1;e.$axios.get(e.baseUrl+"/api/v1/admin/applicants/certs/"+e.curCertId+"?page="+e.currentPage+"&limit="+e.curPageSize+"&type="+e.radio+"&field="+e.searchWord).then(function(a){a.data.success?(e.allListData=a.data.data.list,e.pageTotal=a.data.data.count):e.$message.error(a.data.msg)})},beforeUpload:function(e){this.excelsdata={Authorization:"Bearer "+localStorage.token}},handleSuccessReject:function(e,a,t){var s=this;if(e.success){var r=e.data.excel_save_path;s.$confirm("是否确认退款？").then(function(e){s.getExcelResult("refunded",r)}).catch(function(e){s.$message({message:"已取消",type:"warning"})})}else s.$message.error(e.msg)},handleSuccess:function(e,a,t){if(e.success){var s=e.data.excel_save_path;this.getExcelResult("passed",s)}else this.$message.error(e.msg)},getExcelResult:function(e,a){var t=this,s={action:e,file_path:a};t.$axios.put(t.baseUrl+"/api/v1/admin/applicants/certs/"+t.curCertId,s).then(function(e){e.data.success?(t.$message({message:"操作成功",type:"success"}),t.getData(t.radio,t.currentPage,t.curPageSize)):t.$message.error(e.data.msg)},function(e){t.$message.error(e.msg)})},dateForm:function(e){var a=new Date(1e3*e),t=a.getFullYear(),s=a.getMonth()+1;s=s<10?"0"+s:s;var r=a.getDate();r=r<10?"0"+r:r;var i=a.getHours();i=i<10?"0"+i:i;var l=a.getMinutes(),n=a.getSeconds();return t+"-"+s+"-"+r+" "+i+":"+(l=l<10?"0"+l:l)+":"+(n=n<10?"0"+n:n)}}}),i={render:function(){var e=this,a=e.$createElement,t=e._self._c||a;return t("div",{},[t("div",{staticClass:"rad-group radio-cont"},[t("el-radio-group",{on:{change:e.changeTab},model:{value:e.radio,callback:function(a){e.radio=a},expression:"radio"}},[t("el-radio-button",{attrs:{label:"export"}},[e._v("待导出")]),e._v(" "),t("el-radio-button",{attrs:{label:"verify"}},[e._v("审核中")]),e._v(" "),t("el-radio-button",{attrs:{label:"passed"}},[e._v("已通过")]),e._v(" "),t("el-radio-button",{attrs:{label:"reject"}},[e._v("已拒绝")])],1)],1),e._v(" "),t("div",[t("div",{staticClass:"handle-box"},[t("el-form",{staticClass:"inlineform1",attrs:{inline:!0}},[t("el-form-item",{attrs:{label:""}},[t("el-button",{directives:[{name:"show",rawName:"v-show",value:"export"==e.radio,expression:"radio == 'export'"}],staticClass:"handle-del mr10",attrs:{type:"primary",icon:"plus"},on:{click:e.outsExport}},[e._v("批量导出")]),e._v(" "),t("el-upload",{directives:[{name:"show",rawName:"v-show",value:"verify"==e.radio,expression:"radio == 'verify'"}],staticClass:"unload-mine",attrs:{type:"primary",action:e.uploadUrl,name:"excel",data:e.excelsdata,"show-file-list":!1,"before-upload":e.beforeUpload,"on-success":e.handleSuccess,"file-list":e.fileList}},[t("el-button",{attrs:{size:"small",type:"primary"}},[e._v("批量导入")])],1),e._v(" "),t("el-button",{directives:[{name:"show",rawName:"v-show",value:"verify"==e.radio,expression:"radio == 'verify'"}],staticClass:"handle-del mr10",attrs:{type:"primary",icon:"plus"},on:{click:e.sureReject}},[e._v("确认拒绝")]),e._v(" "),e.rejectDisabled?e._e():t("el-button",{directives:[{name:"show",rawName:"v-show",value:"reject"==e.radio&&e.showTable,expression:"radio == 'reject' && showTable "}],staticClass:"handle-del mr10",attrs:{type:"primary",icon:"plus"},on:{click:e.outsReject}},[e._v("批量导出")]),e._v(" "),e.rejectDisabled?e._e():t("el-upload",{directives:[{name:"show",rawName:"v-show",value:"reject"==e.radio&&e.showTable,expression:"radio == 'reject' && showTable "}],staticClass:"unload-mine",attrs:{type:"primary",action:e.uploadUrl,data:e.excelsdata,name:"excel","show-file-list":!1,"before-upload":e.beforeUpload,"on-success":e.handleSuccessReject,"file-list":e.fileList2}},[e.rejectDisabled?e._e():t("el-button",{attrs:{size:"small",type:"primary"}},[e._v("批量导入")])],1)],1)],1),e._v(" "),t("el-form",{directives:[{name:"show",rawName:"v-show",value:e.showTable,expression:"showTable"}],staticClass:"inlineform2",attrs:{inline:!0}},[t("el-form-item",{attrs:{label:""}},[t("el-input",{staticClass:"handle-input mr10",attrs:{placeholder:"请输入申请人或身份证号查找"},model:{value:e.searchWord,callback:function(a){e.searchWord=a},expression:"searchWord"}})],1),e._v(" "),t("el-form-item",[t("el-button",{attrs:{type:"primary",icon:"search"},on:{click:e.search}},[e._v("查询")])],1)],1)],1),e._v(" "),t("div",{directives:[{name:"show",rawName:"v-show",value:e.showTable,expression:"showTable"}]},[t("el-table",{directives:[{name:"loading",rawName:"v-loading",value:e.loading,expression:"loading"}],key:e.radio,ref:"multipleTable",staticStyle:{width:"100%"},attrs:{data:e.allListData,border:""},on:{"selection-change":e.handleSelectionChange}},[t("el-table-column",{attrs:{type:"index",label:"编号",width:"50"}}),e._v(" "),"verify"==e.radio?t("el-table-column",{attrs:{type:"selection",width:"55"}}):e._e(),e._v(" "),t("el-table-column",{attrs:{prop:"cert_name",label:"申请证书",width:"150"}}),e._v(" "),t("el-table-column",{attrs:{prop:"name",label:"申请人"}}),e._v(" "),t("el-table-column",{attrs:{prop:"personal_id",label:"身份证号",width:"180"}}),e._v(" "),"passed"==e.radio?t("el-table-column",{attrs:{prop:"serial_number",label:"证书编号",width:"150"}}):e._e(),e._v(" "),t("el-table-column",{attrs:{prop:"apply_date",label:"申请时间",width:"180"},scopedSlots:e._u([{key:"default",fn:function(a){return[e._v("\n            "+e._s(e.dateForm(a.row.apply_date))+"\n          ")]}}])}),e._v(" "),t("el-table-column",{attrs:{prop:"pay_amount",label:"支付金额"},scopedSlots:e._u([{key:"default",fn:function(a){return[e._v("\n            "+e._s(a.row.pay_amount+" 元")+"\n          ")]}}])}),e._v(" "),t("el-table-column",{attrs:{prop:"apply_status_msg",label:"申请状态"}}),e._v(" "),"passed"==e.radio?t("el-table-column",{attrs:{prop:"",label:"操作"},scopedSlots:e._u([{key:"default",fn:function(a){return[t("el-button",{attrs:{type:"text",size:"small"},on:{click:function(t){e.clickDetailImg(a.row.cert_id,a.row.cert_name,a.row.wechat_id)}}},[e._v("查看")])]}}])}):e._e(),e._v(" "),"reject"==e.radio?t("el-table-column",{attrs:{prop:"",label:"查看"},scopedSlots:e._u([{key:"default",fn:function(a){return[t("el-button",{attrs:{type:"text",size:"small"},on:{click:function(t){e.clickDetail(a.row.cert_name,a.row.name,a.row.personal_id,a.row.study_date,a.row.pay_amount,a.row.apply_status_msg,a.row.apply_status,a.row.pay_status)}}},[e._v("详情")])]}}])}):e._e()],1),e._v(" "),t("div",{staticClass:"pagination"},[t("el-pagination",{attrs:{"current-page":e.currentPage,layout:"sizes,prev, pager, next","page-sizes":[10,20,30,40],"page-size":e.curPageSize,total:e.pageTotal},on:{"size-change":e.handleSizeChange,"current-change":e.handleCurrentChange,"update:currentPage":function(a){e.currentPage=a}}})],1)],1),e._v(" "),t("div",{directives:[{name:"show",rawName:"v-show",value:!e.showTable,expression:"!showTable"}],staticClass:"detail-cont"},[t("p",[e._v("退款详情")]),e._v(" "),t("ul",{staticClass:"detail-msg"},[t("li",[t("span",[e._v("申请证书：")]),e._v(" "),t("p",[e._v(e._s(e.msg.certname))])]),e._v(" "),t("li",[t("span",[e._v("申请人：")]),e._v(" "),t("p",[e._v(e._s(e.msg.name))])]),e._v(" "),t("li",[t("span",[e._v("身份证号：")]),e._v(" "),t("p",[e._v(e._s(e.msg.cardid))])]),e._v(" "),t("li",[t("span",[e._v("申请时间：")]),e._v(" "),t("p",[e._v(e._s(e.msg.date))])]),e._v(" "),t("li",[t("span",[e._v("支付金额：")]),e._v(" "),t("p",[e._v(e._s(e.msg.price||0))])]),e._v(" "),t("li",[t("span",[e._v("退款进度：")]),e._v(" "),t("p",[e._v(e._s(e.msg.status))])])])])]),e._v(" "),t("div",{directives:[{name:"show",rawName:"v-show",value:"waitexport"==e.radio,expression:"radio == 'waitexport'"}]},[t("w-export")],1),e._v(" "),t("el-dialog",{attrs:{title:e.passedCertName,visible:e.dialogImgVisible,width:"30%","modal-append-to-body":!1},on:{"update:visible":function(a){e.dialogImgVisible=a}}},[t("span",{staticClass:"dialog-img-cont"},[t("img",{attrs:{src:e.passedImg,alt:""}})]),e._v(" "),t("span",{staticClass:"dialog-footer",attrs:{slot:"footer"},slot:"footer"})])],1)},staticRenderFns:[]};var l=t("VU/8")(r,i,!1,function(e){t("7a/k")},"data-v-77b608a0",null);a.default=l.exports},"Eiw/":function(e,a){}});
//# sourceMappingURL=3.66d499a1ea0df52593af.js.map