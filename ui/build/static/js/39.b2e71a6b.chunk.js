(this["webpackJsonp@coreui/coreui-free-react-admin-template"]=this["webpackJsonp@coreui/coreui-free-react-admin-template"]||[]).push([[39],{694:function(e,t,a){"use strict";a.r(t);var n=a(114),r=a.n(n),l=a(170),c=a(63),s=a(64),o=a(107),u=a(387),m=a(386),i=a(1),p=a.n(i),d=a(630),h=a(632),E=a(162),f=a(109),b=function(e){Object(u.a)(a,e);var t=Object(m.a)(a);function a(e){var n;return Object(c.a)(this,a),(n=t.call(this,e)).state={username:"",password:""},n.handleChange=n.handleChange.bind(Object(o.a)(n)),n}return Object(s.a)(a,[{key:"handleChange",value:function(e){if(e.target.name){var t={};t[e.target.name]=e.target.value,this.setState(t),e.preventDefault()}}},{key:"render",value:function(){var e=this;return p.a.createElement(E.a.Consumer,null,(function(t){return p.a.createElement("div",{className:"c-app c-default-layout flex-row align-items-center"},p.a.createElement(d.w,null,p.a.createElement(d.wb,{className:"justify-content-center"},p.a.createElement(d.u,{md:"8"},p.a.createElement(d.j,{className:"p-4"},p.a.createElement(d.k,null,p.a.createElement(d.gb,{show:!!e.state.error,onClose:function(){return e.setState({error:void 0})},color:"danger",size:"sm"},p.a.createElement(d.jb,{closeButton:!0},p.a.createElement(d.kb,null,"Login or password wrong")),p.a.createElement(d.hb,null,e.state.error),p.a.createElement(d.ib,null,p.a.createElement(d.f,{color:"primary",onClick:function(){return e.setState({error:void 0})}},"Ok"))),p.a.createElement(d.J,null,p.a.createElement("h1",null,"Login"),p.a.createElement("p",{className:"text-muted"},"Sign In to your account"),p.a.createElement(d.V,{className:"mb-3"},p.a.createElement(d.X,null,p.a.createElement(d.Y,null,p.a.createElement(h.b,{name:"cil-user"}))),p.a.createElement(d.S,{type:"text",name:"username",placeholder:JSON.stringify(t),autoComplete:"username",value:e.state.username,onChange:e.handleChange})),p.a.createElement(d.V,{className:"mb-4"},p.a.createElement(d.X,null,p.a.createElement(d.Y,null,p.a.createElement(h.b,{name:"cil-lock-locked"}))),p.a.createElement(d.S,{type:"password",name:"password",placeholder:"Password",autoComplete:"current-password",value:e.state.password,onChange:e.handleChange})),p.a.createElement(d.wb,null,p.a.createElement(d.u,{xs:"6"},p.a.createElement(d.f,{color:"primary",className:"px-4",onClick:function(){var a=Object(l.a)(r.a.mark((function a(n){var l;return r.a.wrap((function(a){for(;;)switch(a.prev=a.next){case 0:return a.next=2,t.login(e.state.username,e.state.password,e.props.history);case 2:(l=a.sent)?e.setState({error:f.b[l]}):(e.props.history.push("/"),console.log("success"));case 4:case"end":return a.stop()}}),a)})));return function(e){return a.apply(this,arguments)}}()},"Login"))))))))))}))}}]),a}(p.a.Component);t.default=b}}]);
//# sourceMappingURL=39.b2e71a6b.chunk.js.map