(this["webpackJsonp@coreui/coreui-free-react-admin-template"]=this["webpackJsonp@coreui/coreui-free-react-admin-template"]||[]).push([[8],{631:function(e,t,a){"use strict";a.d(t,"a",(function(){return r}));var n=a(168);function r(e,t){return function(e){if(Array.isArray(e))return e}(e)||function(e,t){if("undefined"!==typeof Symbol&&Symbol.iterator in Object(e)){var a=[],n=!0,r=!1,l=void 0;try{for(var c,s=e[Symbol.iterator]();!(n=(c=s.next()).done)&&(a.push(c.value),!t||a.length!==t);n=!0);}catch(i){r=!0,l=i}finally{try{n||null==s.return||s.return()}finally{if(r)throw l}}return a}}(e,t)||Object(n.a)(e,t)||function(){throw new TypeError("Invalid attempt to destructure non-iterable instance.\nIn order to be iterable, non-array objects must have a [Symbol.iterator]() method.")}()}},686:function(e,t,a){"use strict";a.r(t);var n=a(114),r=a.n(n),l=a(170),c=a(631),s=a(1),i=a.n(s),o=a(630),u=a(632),p=function(e){return i.a.createElement(o.u,{xl:6},i.a.createElement(o.j,null,i.a.createElement(o.n,null,i.a.createElement(u.a,{name:"cil-applications"})," ",e.name," - ",i.a.createElement("code",null,e.hash),i.a.createElement("div",{className:"float-right"},"/"+e.mountPoint)),i.a.createElement(o.k,null,!!e.description&&i.a.createElement("p",{class:"lead"},i.a.createElement("em",null,e.description)),i.a.createElement("h5",null,"Available Roles"),i.a.createElement(o.y,{striped:!0,items:e.availableRoles,fields:["name","description","autoAssign"],size:"sm"}),i.a.createElement("h5",null,"Access Policies"),i.a.createElement(o.y,{striped:!0,items:e.accessPolicies,fields:["roles","resources","effect","actions"],size:"sm"}))))};t.default=function(){var e=Object(s.useState)({data:[]}),t=Object(c.a)(e,2),a=t[0],n=t[1];return Object(s.useEffect)((function(){var e=!1;function t(){return(t=Object(l.a)(r.a.mark((function t(){var a;return r.a.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return t.next=2,fetch("http://localhost/admin/api/v0/apps",{method:"GET"});case 2:if(200!==(a=t.sent).status){t.next=10;break}if(e){t.next=10;break}return t.t0=n,t.next=8,a.json();case 8:t.t1=t.sent,(0,t.t0)(t.t1);case 10:case"end":return t.stop()}}),t)})))).apply(this,arguments)}return function(){t.apply(this,arguments)}(),function(){e=!0}}),[]),i.a.createElement(i.a.Fragment,null,i.a.createElement(o.wb,null,i.a.createElement(o.u,null,i.a.createElement(o.j,null,i.a.createElement(o.k,null,i.a.createElement("p",null,"The following cypherapps have been installed using the ",i.a.createElement("code",null,"cam")," tool. They cannot (yet) be edited using the cyphernode admin. It's just a quick overview of all the roles a cypherapp provides and how access to different resources is handled."),i.a.createElement("p",null,i.a.createElement("h5",null,"Available Roles"),"The available roles describe the roles present in the respective cypherapp. Please note that all roles are local and for example the ",i.a.createElement("code",null,"admin")," role in one cypherapp is not the same as the ",i.a.createElement("code",null,"admin")," role in another cypherapp. If a role has the ",i.a.createElement("code",null,"autoAssign")," flag set, it means that every user will automatically receive that role. This is useful, when you want to install a new cypherapp and give the existing users access to the regular user space of that cypherapp."),i.a.createElement("p",null,i.a.createElement("h5",null,"Access Policies"),"Access policies allow or deny access of certain roles within the cypherapp to a certain resource relative to the mount point of the cypherapp. All resources are described using regular expressions. The actions are equivalent to the available HTTP methods."))))),i.a.createElement(o.wb,null," ",a.data.map((function(e){return i.a.createElement(p,e)}))," "))}}}]);
//# sourceMappingURL=8.1c70f3d0.chunk.js.map