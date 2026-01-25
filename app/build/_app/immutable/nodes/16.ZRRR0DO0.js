import"../chunks/DsnmJJEf.js";import{a as Ye}from"../chunks/Hdby3CWe.js";import{f as $,p as Be,a as xe,s as T,d as He,e as qe,b as c,h as r,i as e,c as s,j as w,g as De,n as U,r as a,t as te}from"../chunks/odI--4Kw.js";import{s as j}from"../chunks/BNMK9zrr.js";import{l as Me,s as ze,i as S}from"../chunks/COaB7F5h.js";import{e as Ne}from"../chunks/BohEjZuL.js";import{c as B,a as l,f as p}from"../chunks/CCInJaqd.js";import{h as We}from"../chunks/hUjpWz_l.js";import{a as Ge,s as Je,C as D,M as Fe,B as M,q as Ve,I as ae,e as Xe,g as he,b as _e,f as ye,A as se}from"../chunks/DVgw2y5W.js";import{L as Qe}from"../chunks/D716f0vp.js";import{a2 as Ze,a3 as et}from"../chunks/cAHYBDdg.js";import{P as tt}from"../chunks/DODF9Coh.js";import"../chunks/B9hjHVe6.js";import{K as at}from"../chunks/jvKZ-wEy.js";import{T as st}from"../chunks/sVcY4UvN.js";function ot(z,H){const N=Me(H,["children","$$slots","$$events","$$legacy"]);const v=[["path",{d:"M11 21.73a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73z"}],["path",{d:"M12 22V12"}],["polyline",{points:"3.29 7 12 12 20.71 7"}],["path",{d:"m7.5 4.27 9 5.15"}]];Ge(z,ze({name:"package"},()=>N,{get iconNode(){return v},children:(O,R)=>{var C=B(),E=$(C);Je(E,H,"default",{}),l(O,C)},$$slots:{default:!0}}))}var rt=p("<p> </p>"),lt=p("<p> </p>"),nt=p("<!> Create Key",1),ct=p('<div class="flex justify-center py-8"><!></div>'),it=p('<div class="mt-6 text-center py-8 text-text-muted"><!> <p class="mt-2">No API keys yet</p> <p class="text-sm">Create an API key to use the REST API or SMTP</p></div>'),dt=p('<tr class="border-b border-border/50"><td class="py-3 font-medium"> </td><td class="py-3 font-mono text-text-muted"> </td><td class="py-3 text-text-muted"> </td><td class="py-3 text-text-muted"> </td><td class="py-3 text-right"><!></td></tr>'),pt=p('<div class="mt-6 overflow-x-auto"><table class="w-full text-sm"><thead><tr class="border-b border-border"><th class="text-left py-2 font-medium text-text-muted">Name</th><th class="text-left py-2 font-medium text-text-muted">Key</th><th class="text-left py-2 font-medium text-text-muted">Last Used</th><th class="text-left py-2 font-medium text-text-muted">Created</th><th class="text-right py-2 font-medium text-text-muted">Actions</th></tr></thead><tbody></tbody></table></div>'),mt=p('<div class="flex items-center justify-between"><div><h2 class="text-lg font-medium text-text mb-1">API Keys</h2> <p class="text-sm text-text-muted">Manage API keys for REST API and SMTP authentication</p></div> <!></div> <!>',1),ut=p('<div class="flex items-center gap-3 mb-4"><!> <div><h2 class="text-lg font-medium text-text">REST API</h2> <p class="text-sm text-text-muted">Connect to Outlet programmatically via the REST API</p></div></div> <div class="space-y-4"><div><label class="form-label">Base URL</label> <div class="flex gap-2"><!> <!></div></div> <div><label class="form-label">Authentication</label> <p class="text-sm text-text-muted mb-2">Include your API key in the Authorization header:</p> <pre class="bg-surface-tertiary p-3 rounded text-sm overflow-x-auto font-mono">Authorization: Bearer YOUR_API_KEY</pre></div></div>',1),vt=p('<h2 class="text-lg font-medium text-text mb-4">Code Examples</h2> <div class="mb-4"><!></div> <!>',1),ft=p('<div class="flex items-center gap-3 mb-4"><!> <div><h2 class="text-lg font-medium text-text">Official SDKs</h2> <p class="text-sm text-text-muted">Pre-built client libraries for popular languages</p></div></div> <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4"><div class="p-4 bg-surface-secondary rounded-lg"><div class="flex items-center gap-2 mb-2"><span class="text-lg font-medium text-text">TypeScript</span></div> <p class="text-sm text-text-muted mb-3">Full TypeScript support with types</p> <code class="block text-xs bg-surface-tertiary px-2 py-1 rounded font-mono text-text-muted">npm install @outlet/sdk</code></div> <div class="p-4 bg-surface-secondary rounded-lg"><div class="flex items-center gap-2 mb-2"><span class="text-lg font-medium text-text">Python</span></div> <p class="text-sm text-text-muted mb-3">Python 3.8+ support</p> <code class="block text-xs bg-surface-tertiary px-2 py-1 rounded font-mono text-text-muted">pip install outlet-sdk</code></div> <div class="p-4 bg-surface-secondary rounded-lg"><div class="flex items-center gap-2 mb-2"><span class="text-lg font-medium text-text">Go</span></div> <p class="text-sm text-text-muted mb-3">Go 1.20+ support</p> <code class="block text-xs bg-surface-tertiary px-2 py-1 rounded font-mono text-text-muted">go get github.com/outlet/sdk-go</code></div> <div class="p-4 bg-surface-secondary rounded-lg"><div class="flex items-center gap-2 mb-2"><span class="text-lg font-medium text-text">PHP</span></div> <p class="text-sm text-text-muted mb-3">PHP 8.0+ support</p> <code class="block text-xs bg-surface-tertiary px-2 py-1 rounded font-mono text-text-muted">composer require outlet/sdk</code></div></div> <div class="mt-6 pt-4 border-t border-border"><h3 class="font-medium text-text mb-3">SDK Quick Start (TypeScript)</h3> <!></div>',1),xt=p("<p>This is the only time you'll see this key. Copy it now and store it securely.</p>"),ht=p('<div class="space-y-4"><!> <div><label class="form-label">API Key</label> <div class="flex gap-2"><!> <!></div></div></div>'),_t=p('<div class="space-y-4"><div><label for="key-name" class="form-label">Key Name</label> <!> <p class="mt-1 text-xs text-text-muted">A descriptive name to identify this key</p></div></div>'),yt=p('<div class="space-y-6"><!> <!> <!> <!> <!> <!></div> <!>',1);function jt(z,H){Be(H,!0);let N=xe(window.location.origin.replace(":5173",":8888")),v=w(()=>`${N}/sdk/v1`),O=T(!0),R=T(""),C=T(""),E=T(""),W=T(xe([])),q=T(!1),G=T(""),L=T("");const ge=[{id:"curl",label:"cURL"},{id:"nodejs",label:"Node.js"},{id:"python",label:"Python"},{id:"go",label:"Go"},{id:"php",label:"PHP"}];let J=T("curl");const be=w(()=>`# List contacts
curl -X GET "${e(v)}/contacts" \\
  -H "Authorization: Bearer YOUR_API_KEY" \\
  -H "Content-Type: application/json"

# Create a contact
curl -X POST "${e(v)}/contacts" \\
  -H "Authorization: Bearer YOUR_API_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{
    "email": "user@example.com",
    "name": "John Doe",
    "tags": ["customer", "newsletter"]
  }'

# Send transactional email
curl -X POST "${e(v)}/emails/send" \\
  -H "Authorization: Bearer YOUR_API_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{
    "to": "recipient@example.com",
    "subject": "Welcome!",
    "html": "<h1>Hello</h1><p>Welcome to our service.</p>"
  }'`),$e=w(()=>`import axios from 'axios';

const api = axios.create({
  baseURL: '${e(v)}',
  headers: {
    'Authorization': 'Bearer YOUR_API_KEY',
    'Content-Type': 'application/json'
  }
});

// List contacts
const contacts = await api.get('/contacts');
console.log(contacts.data);

// Create a contact
const newContact = await api.post('/contacts', {
  email: 'user@example.com',
  name: 'John Doe',
  tags: ['customer', 'newsletter']
});

// Send transactional email
await api.post('/emails/send', {
  to: 'recipient@example.com',
  subject: 'Welcome!',
  html: '<h1>Hello</h1><p>Welcome to our service.</p>'
});`),Pe=w(()=>`import requests

API_KEY = 'YOUR_API_KEY'
BASE_URL = '${e(v)}'
headers = {
    'Authorization': f'Bearer {API_KEY}',
    'Content-Type': 'application/json'
}

# List contacts
response = requests.get(f'{BASE_URL}/contacts', headers=headers)
contacts = response.json()

# Create a contact
new_contact = requests.post(f'{BASE_URL}/contacts', headers=headers, json={
    'email': 'user@example.com',
    'name': 'John Doe',
    'tags': ['customer', 'newsletter']
})

# Send transactional email
requests.post(f'{BASE_URL}/emails/send', headers=headers, json={
    'to': 'recipient@example.com',
    'subject': 'Welcome!',
    'html': '<h1>Hello</h1><p>Welcome to our service.</p>'
})`),Ae=w(()=>`package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

const (
    apiKey  = "YOUR_API_KEY"
    baseURL = "${e(v)}"
)

func main() {
    client := &http.Client{}

    // Create contact
    contact := map[string]interface{}{
        "email": "user@example.com",
        "name":  "John Doe",
        "tags":  []string{"customer", "newsletter"},
    }
    body, _ := json.Marshal(contact)

    req, _ := http.NewRequest("POST", baseURL+"/contacts", bytes.NewBuffer(body))
    req.Header.Set("Authorization", "Bearer "+apiKey)
    req.Header.Set("Content-Type", "application/json")

    resp, _ := client.Do(req)
    defer resp.Body.Close()
}`),Te=w(()=>`<?php

$apiKey = 'YOUR_API_KEY';
$baseUrl = '${e(v)}';

// Helper function for API requests
function apiRequest($method, $endpoint, $data = null) {
    global $apiKey, $baseUrl;

    $ch = curl_init($baseUrl . $endpoint);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_CUSTOMREQUEST, $method);
    curl_setopt($ch, CURLOPT_HTTPHEADER, [
        'Authorization: Bearer ' . $apiKey,
        'Content-Type: application/json'
    ]);

    if ($data) {
        curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));
    }

    $response = curl_exec($ch);
    curl_close($ch);
    return json_decode($response, true);
}

// List contacts
$contacts = apiRequest('GET', '/contacts');

// Create a contact
$newContact = apiRequest('POST', '/contacts', [
    'email' => 'user@example.com',
    'name' => 'John Doe',
    'tags' => ['customer', 'newsletter']
]);

// Send transactional email
apiRequest('POST', '/emails/send', [
    'to' => 'recipient@example.com',
    'subject' => 'Welcome!',
    'html' => '<h1>Hello</h1><p>Welcome to our service.</p>'
]);`);let we=w(()=>({curl:{code:e(be),language:"bash"},nodejs:{code:e($e),language:"javascript"},python:{code:e(Pe),language:"python"},go:{code:e(Ae),language:"go"},php:{code:e(Te),language:"php"}})),oe=w(()=>e(we)[e(J)]);Ye(async()=>{await re()});async function re(){c(O,!0),c(R,"");try{const t=await Ze();c(W,t.keys||[],!0)}catch{console.log("No API keys yet")}finally{c(O,!1)}}async function Ce(t){if(confirm("Are you sure you want to revoke this API key? This cannot be undone."))try{await et({},t),await re(),c(C,"API key revoked")}catch(P){c(R,P.message||"Failed to revoke API key",!0)}}function le(t){navigator.clipboard.writeText(t),c(E,t,!0),setTimeout(()=>c(E,""),2e3)}function ke(){c(q,!1),c(G,""),c(L,"")}function ne(t){return t?new Date(t).toLocaleDateString("en-US",{month:"short",day:"numeric",year:"numeric"}):"Never"}var ce=yt();We("1kyawqp",t=>{qe(()=>{De.title="API - Settings"})});var F=$(ce),ie=s(F);{var Ie=t=>{se(t,{type:"error",title:"Error",onclose:()=>c(R,""),children:(P,i)=>{var o=rt(),f=s(o,!0);a(o),te(()=>j(f,e(R))),l(P,o)},$$slots:{default:!0}})};S(ie,t=>{e(R)&&t(Ie)})}var de=r(ie,2);{var Se=t=>{se(t,{type:"success",title:"Success",onclose:()=>c(C,""),children:(P,i)=>{var o=lt(),f=s(o,!0);a(o),te(()=>j(f,e(C))),l(P,o)},$$slots:{default:!0}})};S(de,t=>{e(C)&&t(Se)})}var pe=r(de,2);D(pe,{children:(t,P)=>{var i=mt(),o=$(i),f=r(s(o),2);M(f,{type:"primary",size:"sm",onclick:()=>c(q,!0),children:(n,u)=>{var x=nt(),k=$(x);tt(k,{class:"mr-1.5 h-4 w-4"}),U(),l(n,x)},$$slots:{default:!0}}),a(o);var _=r(o,2);{var d=n=>{var u=ct(),x=s(u);Qe(x,{}),a(u),l(n,u)},m=n=>{var u=B(),x=$(u);{var k=h=>{var y=it(),g=s(y);at(g,{class:"mx-auto h-12 w-12 opacity-50"}),U(4),a(y),l(h,y)},K=h=>{var y=pt(),g=s(y),b=r(s(g));Ne(b,21,()=>e(W),Y=>Y.id,(Y,I)=>{var A=dt(),V=s(A),Ke=s(V,!0);a(V);var X=r(V),Ue=s(X);a(X);var Q=r(X),je=s(Q,!0);a(Q);var Z=r(Q),Oe=s(Z,!0);a(Z);var ve=r(Z),Le=s(ve);M(Le,{type:"danger",size:"icon",onclick:()=>Ce(e(I).id),children:(ee,fe)=>{st(ee,{class:"h-4 w-4"})},$$slots:{default:!0}}),a(ve),a(A),te((ee,fe)=>{j(Ke,e(I).name),j(Ue,`${e(I).key_prefix??""}...`),j(je,ee),j(Oe,fe)},[()=>ne(e(I).last_used),()=>ne(e(I).created_at)]),l(Y,A)}),a(b),a(g),a(y),l(h,y)};S(x,h=>{e(W).length===0?h(k):h(K,!1)},!0)}l(n,u)};S(_,n=>{e(O)?n(d):n(m,!1)})}l(t,i)},$$slots:{default:!0}});var me=r(pe,2);D(me,{children:(t,P)=>{var i=ut(),o=$(i),f=s(o);Ve(f,{class:"h-5 w-5 text-primary"}),U(2),a(o);var _=r(o,2),d=s(_),m=r(s(d),2),n=s(m);ae(n,{type:"text",get value(){return e(v)},readonly:!0,class:"font-mono text-sm"});var u=r(n,2);M(u,{type:"secondary",onclick:()=>le(e(v)),children:(x,k)=>{var K=B(),h=$(K);{var y=b=>{_e(b,{class:"h-4 w-4 text-green-500"})},g=b=>{ye(b,{class:"h-4 w-4"})};S(h,b=>{e(E)===e(v)?b(y):b(g,!1)})}l(x,K)},$$slots:{default:!0}}),a(m),a(d),U(2),a(_),l(t,i)},$$slots:{default:!0}});var ue=r(me,2);D(ue,{children:(t,P)=>{var i=vt(),o=r($(i),2),f=s(o);Xe(f,{get tabs(){return ge},variant:"pills",get activeTab(){return e(J)},set activeTab(d){c(J,d,!0)}}),a(o);var _=r(o,2);he(_,{get code(){return e(oe).code},get language(){return e(oe).language}}),l(t,i)},$$slots:{default:!0}});var Re=r(ue,2);D(Re,{children:(t,P)=>{var i=ft(),o=$(i),f=s(o);ot(f,{class:"h-5 w-5 text-primary"}),U(2),a(o);var _=r(o,4),d=r(s(_),2);{let m=w(()=>`import { Outlet } from '@outlet/sdk';

const outlet = new Outlet('YOUR_API_KEY', '${e(v).replace("/sdk/v1","")}');

// Send a transactional email
const result = await outlet.emails.sendEmail({
  to: 'user@example.com',
  subject: 'Welcome!',
  html_body: '<h1>Welcome to our platform</h1>',
});

console.log('Message ID:', result.message_id);`);he(d,{get code(){return e(m)},language:"typescript"})}a(_),l(t,i)},$$slots:{default:!0}}),a(F);var Ee=r(F,2);Fe(Ee,{title:"Create API Key",onclose:ke,get show(){return e(q)},set show(t){c(q,t,!0)},children:(t,P)=>{var i=B(),o=$(i);{var f=d=>{var m=ht(),n=s(m);se(n,{type:"warning",title:"Save your API key",children:(h,y)=>{var g=xt();l(h,g)},$$slots:{default:!0}});var u=r(n,2),x=r(s(u),2),k=s(x);ae(k,{type:"text",get value(){return e(L)},readonly:!0,class:"font-mono text-sm"});var K=r(k,2);M(K,{type:"secondary",onclick:()=>le(e(L)),children:(h,y)=>{var g=B(),b=$(g);{var Y=A=>{_e(A,{class:"h-4 w-4 text-green-500"})},I=A=>{ye(A,{class:"h-4 w-4"})};S(b,A=>{e(E)===e(L)?A(Y):A(I,!1)})}l(h,g)},$$slots:{default:!0}}),a(x),a(u),a(m),l(d,m)},_=d=>{var m=_t(),n=s(m),u=r(s(n),2);ae(u,{type:"text",id:"key-name",placeholder:"e.g., Production API, Development",get value(){return e(G)},set value(x){c(G,x,!0)}}),U(2),a(n),a(m),l(d,m)};S(o,d=>{e(L)?d(f):d(_,!1)})}l(t,i)},$$slots:{default:!0}}),l(z,ce),He()}export{jt as component};
