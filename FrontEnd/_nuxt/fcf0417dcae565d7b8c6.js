(window.webpackJsonp=window.webpackJsonp||[]).push([[2],{224:function(t,e,r){"use strict";r.r(e);var n={data:function(){return{form:{url:"",response:[]},show:!0,showActivity:!1}},methods:{onSubmit:function(t){var e=this;t.preventDefault(),this.showActivity=!0,this.$axios.get("http://192.168.1.61:8080/api/checkServer/"+encodeURIComponent(this.form.url)).then((function(t){e.form.response={items:[t.data]},e.showActivity=!1}))},onReset:function(t){var e=this;t.preventDefault(),this.form.url="",this.show=!1,this.$nextTick((function(){e.show=!0}))},actionHistory:function(t){var e=this;t.preventDefault(),this.showActivity=!0,this.$axios.get("http://192.168.1.61:8080/api/checkServerHistory").then((function(t){e.form.response=t.data,e.showActivity=!1}))}}},o=r(48),component=Object(o.a)(n,(function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("div",{staticClass:"container"},[r("h1",[t._v("Fync Url Checker")]),t._v(" "),t.show?r("b-form",{on:{submit:t.onSubmit,reset:t.onReset}},[r("b-form-group",{attrs:{id:"input-group-1",label:"Url To Check:","label-for":"input-1",description:"Write the url to check"}},[r("b-form-input",{attrs:{id:"input-1",type:"text",required:"",placeholder:"Enter url"},model:{value:t.form.url,callback:function(e){t.$set(t.form,"url",e)},expression:"form.url"}})],1),t._v(" "),r("b-button",{attrs:{type:"submit",variant:"primary"}},[t._v("Submit")]),t._v(" "),r("b-button",{attrs:{type:"reset",variant:"danger"}},[t._v("Reset")]),t._v(" "),r("b-button",{attrs:{type:"button",variant:"secondary"},on:{click:t.actionHistory}},[t._v("History")]),t._v(" "),t.showActivity?r("b-spinner",{attrs:{type:"grow",label:"Loading..."}}):t._e()],1):t._e(),t._v(" "),t._l(t.form.response.items,(function(e){return[r("b-card",{key:e,staticClass:"mt-4",attrs:{title:e.url,"header-tag":"header"},scopedSlots:t._u([{key:"header",fn:function(){return[r("h6",{staticClass:"mb-0"},[r("b-avatar",{attrs:{variant:"light",src:e.logo,icon:"star-fill",size:"4em"}}),t._v("\n          "+t._s(e.url)+" "),e.title?r("span",[t._v("- "+t._s(e.title))]):t._e(),t._v(" "),r("b-badge",{attrs:{variant:"primary"}},[t._v(t._s(e.ssl_grade)+" ")]),t._v(" "),r("b-badge",[t._v(t._s(e.previous_ssl_grade)+" ")]),t._v(" "),e.is_down?r("b-badge",{attrs:{variant:"danger"}},[t._v("Is down")]):r("b-badge",{attrs:{variant:"success"}},[t._v("Ok")]),t._v(" "),e.servers_changed?r("b-badge",{attrs:{variant:"warning"}},[t._v("Servers changed")]):t._e()],1)]},proxy:!0}],null,!0)},[t._v(" "),r("b-table",{attrs:{striped:"",hover:"",items:e.servers}})],1)]})),t._v(" "),r("b-card",{staticClass:"mt-3",attrs:{header:"Form Data Result"}},[r("pre",{staticClass:"m-0"},[t._v(t._s(t.form))])])],2)}),[],!1,null,null,null);e.default=component.exports}}]);