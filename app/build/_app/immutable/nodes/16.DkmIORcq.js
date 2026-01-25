import"../chunks/DsnmJJEf.js";import{a as De}from"../chunks/Hdby3CWe.js";import{f as P,p as qe,a as _e,s as T,d as Me,e as ze,b as n,h as r,i as e,c as o,j as k,g as Ne,n as I,r as s,t as W}from"../chunks/odI--4Kw.js";import{s as j}from"../chunks/BNMK9zrr.js";import{l as We,s as Fe,i as R}from"../chunks/COaB7F5h.js";import{e as Ge}from"../chunks/BohEjZuL.js";import{c as q,a as l,f as v,t as re}from"../chunks/CCInJaqd.js";import{h as Je}from"../chunks/hUjpWz_l.js";import{a as Ve,s as Xe,C as F,M as Qe,B as O,r as Ze,I as le,e as et,n as ge,b as be,f as $e,A as ne}from"../chunks/C9m1dhUm.js";import{L as tt}from"../chunks/D716f0vp.js";import{a3 as at,a4 as st,a5 as ot}from"../chunks/BYFcJoO7.js";import{P as rt}from"../chunks/Cmj1ijxE.js";import"../chunks/B9hjHVe6.js";import{K as lt}from"../chunks/q1YT1ipT.js";import{T as nt}from"../chunks/BNdV4B7Z.js";function ct(G,M){const J=We(M,["children","$$slots","$$events","$$legacy"]);const h=[["path",{d:"M11 21.73a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73z"}],["path",{d:"M12 22V12"}],["polyline",{points:"3.29 7 12 12 20.71 7"}],["path",{d:"m7.5 4.27 9 5.15"}]];Ve(G,Fe({name:"package"},()=>J,{get iconNode(){return h},children:(Y,B)=>{var A=q(),K=P(A);Xe(K,M,"default",{}),l(Y,A)},$$slots:{default:!0}}))}var it=v("<p> </p>"),dt=v("<p> </p>"),pt=v("<!> Create Key",1),ut=v('<div class="flex justify-center py-8"><!></div>'),mt=v('<div class="mt-6 text-center py-8 text-text-muted"><!> <p class="mt-2">No API keys yet</p> <p class="text-sm">Create an API key to use the REST API or SMTP</p></div>'),vt=v('<tr class="border-b border-border/50"><td class="py-3 font-medium"> </td><td class="py-3 font-mono text-text-muted"> </td><td class="py-3 text-text-muted"> </td><td class="py-3 text-text-muted"> </td><td class="py-3 text-right"><!></td></tr>'),ft=v('<div class="mt-6 overflow-x-auto"><table class="w-full text-sm"><thead><tr class="border-b border-border"><th class="text-left py-2 font-medium text-text-muted">Name</th><th class="text-left py-2 font-medium text-text-muted">Key</th><th class="text-left py-2 font-medium text-text-muted">Last Used</th><th class="text-left py-2 font-medium text-text-muted">Created</th><th class="text-right py-2 font-medium text-text-muted">Actions</th></tr></thead><tbody></tbody></table></div>'),xt=v('<div class="flex items-center justify-between"><div><h2 class="text-lg font-medium text-text mb-1">API Keys</h2> <p class="text-sm text-text-muted">Manage API keys for REST API and SMTP authentication</p></div> <!></div> <!>',1),yt=v('<div class="flex items-center gap-3 mb-4"><!> <div><h2 class="text-lg font-medium text-text">REST API</h2> <p class="text-sm text-text-muted">Connect to Outlet programmatically via the REST API</p></div></div> <div class="space-y-4"><div><label class="form-label">Base URL</label> <div class="flex gap-2"><!> <!></div></div> <div><label class="form-label">Authentication</label> <p class="text-sm text-text-muted mb-2">Include your API key in the Authorization header:</p> <pre class="bg-surface-tertiary p-3 rounded text-sm overflow-x-auto font-mono">Authorization: Bearer YOUR_API_KEY</pre></div></div>',1),ht=v('<h2 class="text-lg font-medium text-text mb-4">Code Examples</h2> <div class="mb-4"><!></div> <!>',1),_t=v('<div class="flex items-center gap-3 mb-4"><!> <div><h2 class="text-lg font-medium text-text">Official SDKs</h2> <p class="text-sm text-text-muted">Pre-built client libraries for popular languages</p></div></div> <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4"><div class="p-4 bg-surface-secondary rounded-lg"><div class="flex items-center gap-2 mb-2"><span class="text-lg font-medium text-text">TypeScript</span></div> <p class="text-sm text-text-muted mb-3">Full TypeScript support with types</p> <code class="block text-xs bg-surface-tertiary px-2 py-1 rounded font-mono text-text-muted">npm install @outlet/sdk</code></div> <div class="p-4 bg-surface-secondary rounded-lg"><div class="flex items-center gap-2 mb-2"><span class="text-lg font-medium text-text">Python</span></div> <p class="text-sm text-text-muted mb-3">Python 3.8+ support</p> <code class="block text-xs bg-surface-tertiary px-2 py-1 rounded font-mono text-text-muted">pip install outlet-sdk</code></div> <div class="p-4 bg-surface-secondary rounded-lg"><div class="flex items-center gap-2 mb-2"><span class="text-lg font-medium text-text">Go</span></div> <p class="text-sm text-text-muted mb-3">Go 1.20+ support</p> <code class="block text-xs bg-surface-tertiary px-2 py-1 rounded font-mono text-text-muted">go get github.com/outlet/sdk-go</code></div> <div class="p-4 bg-surface-secondary rounded-lg"><div class="flex items-center gap-2 mb-2"><span class="text-lg font-medium text-text">PHP</span></div> <p class="text-sm text-text-muted mb-3">PHP 8.0+ support</p> <code class="block text-xs bg-surface-tertiary px-2 py-1 rounded font-mono text-text-muted">composer require outlet/sdk</code></div></div> <div class="mt-6 pt-4 border-t border-border"><h3 class="font-medium text-text mb-3">SDK Quick Start (TypeScript)</h3> <!></div>',1),gt=v("<!> <!>",1),bt=v('<div class="flex justify-end gap-3"><!></div>'),$t=v("<p>This is the only time you'll see this key. Copy it now and store it securely.</p>"),Pt=v('<div class="space-y-4"><!> <div><label class="form-label">API Key</label> <div class="flex gap-2"><!> <!></div></div></div>'),At=v('<div class="space-y-4"><div><label for="key-name" class="form-label">Key Name</label> <!> <p class="mt-1 text-xs text-text-muted">A descriptive name to identify this key</p></div></div>'),wt=v('<div class="space-y-6"><!> <!> <!> <!> <!> <!></div> <!>',1);function Dt(G,M){qe(M,!0);let J=_e(window.location.origin.replace(":5173",":8888")),h=k(()=>`${J}/sdk/v1`),Y=T(!0),B=T(!1),A=T(""),K=T(""),z=T(""),V=T(_e([])),N=T(!1),L=T(""),E=T("");const Pe=[{id:"curl",label:"cURL"},{id:"nodejs",label:"Node.js"},{id:"python",label:"Python"},{id:"go",label:"Go"},{id:"php",label:"PHP"}];let X=T("curl");const Ae=k(()=>`# List contacts
curl -X GET "${e(h)}/contacts" \\
  -H "Authorization: Bearer YOUR_API_KEY" \\
  -H "Content-Type: application/json"

# Create a contact
curl -X POST "${e(h)}/contacts" \\
  -H "Authorization: Bearer YOUR_API_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{
    "email": "user@example.com",
    "name": "John Doe",
    "tags": ["customer", "newsletter"]
  }'

# Send transactional email
curl -X POST "${e(h)}/emails/send" \\
  -H "Authorization: Bearer YOUR_API_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{
    "to": "recipient@example.com",
    "subject": "Welcome!",
    "html": "<h1>Hello</h1><p>Welcome to our service.</p>"
  }'`),we=k(()=>`import axios from 'axios';

const api = axios.create({
  baseURL: '${e(h)}',
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
});`),Ce=k(()=>`import requests

API_KEY = 'YOUR_API_KEY'
BASE_URL = '${e(h)}'
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
})`),Te=k(()=>`package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

const (
    apiKey  = "YOUR_API_KEY"
    baseURL = "${e(h)}"
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
}`),ke=k(()=>`<?php

$apiKey = 'YOUR_API_KEY';
$baseUrl = '${e(h)}';

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
]);`);let Ie=k(()=>({curl:{code:e(Ae),language:"bash"},nodejs:{code:e(we),language:"javascript"},python:{code:e(Ce),language:"python"},go:{code:e(Te),language:"go"},php:{code:e(ke),language:"php"}})),ce=k(()=>e(Ie)[e(X)]);De(async()=>{await Q()});async function Q(){n(Y,!0),n(A,"");try{const t=await at();n(V,t.keys||[],!0)}catch{console.log("No API keys yet")}finally{n(Y,!1)}}async function Se(){if(!e(L).trim()){n(A,"Please enter a name for the API key");return}n(B,!0),n(A,"");try{const t=await ot({name:e(L).trim()});n(E,t.key||"",!0),await Q(),n(K,"API key created successfully. Copy it now - you won't be able to see it again!")}catch(t){n(A,t.message||"Failed to create API key",!0)}finally{n(B,!1)}}async function Re(t){if(confirm("Are you sure you want to revoke this API key? This cannot be undone."))try{await st({},t),await Q(),n(K,"API key revoked")}catch(f){n(A,f.message||"Failed to revoke API key",!0)}}function ie(t){navigator.clipboard.writeText(t),n(z,t,!0),setTimeout(()=>n(z,""),2e3)}function Z(){n(N,!1),n(L,""),n(E,"")}function de(t){return t?new Date(t).toLocaleDateString("en-US",{month:"short",day:"numeric",year:"numeric"}):"Never"}var pe=wt();Je("1kyawqp",t=>{ze(()=>{Ne.title="API - Settings"})});var ee=P(pe),ue=o(ee);{var Ke=t=>{ne(t,{type:"error",title:"Error",onclose:()=>n(A,""),children:(f,d)=>{var a=it(),x=o(a,!0);s(a),W(()=>j(x,e(A))),l(f,a)},$$slots:{default:!0}})};R(ue,t=>{e(A)&&t(Ke)})}var me=r(ue,2);{var Ee=t=>{ne(t,{type:"success",title:"Success",onclose:()=>n(K,""),children:(f,d)=>{var a=dt(),x=o(a,!0);s(a),W(()=>j(x,e(K))),l(f,a)},$$slots:{default:!0}})};R(me,t=>{e(K)&&t(Ee)})}var ve=r(me,2);F(ve,{children:(t,f)=>{var d=xt(),a=P(d),x=r(o(a),2);O(x,{type:"primary",size:"sm",onclick:()=>n(N,!0),children:(c,i)=>{var m=pt(),$=P(m);rt($,{class:"mr-1.5 h-4 w-4"}),I(),l(c,m)},$$slots:{default:!0}}),s(a);var _=r(a,2);{var u=c=>{var i=ut(),m=o(i);tt(m,{}),s(i),l(c,i)},p=c=>{var i=q(),m=P(i);{var $=g=>{var b=mt(),S=o(b);lt(S,{class:"mx-auto h-12 w-12 opacity-50"}),I(4),s(b),l(g,b)},w=g=>{var b=ft(),S=o(b),y=r(o(S));Ge(y,21,()=>e(V),H=>H.id,(H,U)=>{var D=vt(),C=o(D),Oe=o(C,!0);s(C);var te=r(C),Le=o(te);s(te);var ae=r(te),Ye=o(ae,!0);s(ae);var se=r(ae),Be=o(se,!0);s(se);var ye=r(se),He=o(ye);O(He,{type:"danger",size:"icon",onclick:()=>Re(e(U).id),children:(oe,he)=>{nt(oe,{class:"h-4 w-4"})},$$slots:{default:!0}}),s(ye),s(D),W((oe,he)=>{j(Oe,e(U).name),j(Le,`${e(U).key_prefix??""}...`),j(Ye,oe),j(Be,he)},[()=>de(e(U).last_used),()=>de(e(U).created_at)]),l(H,D)}),s(y),s(S),s(b),l(g,b)};R(m,g=>{e(V).length===0?g($):g(w,!1)},!0)}l(c,i)};R(_,c=>{e(Y)?c(u):c(p,!1)})}l(t,d)},$$slots:{default:!0}});var fe=r(ve,2);F(fe,{children:(t,f)=>{var d=yt(),a=P(d),x=o(a);Ze(x,{class:"h-5 w-5 text-primary"}),I(2),s(a);var _=r(a,2),u=o(_),p=r(o(u),2),c=o(p);le(c,{type:"text",get value(){return e(h)},readonly:!0,class:"font-mono text-sm"});var i=r(c,2);O(i,{type:"secondary",onclick:()=>ie(e(h)),children:(m,$)=>{var w=q(),g=P(w);{var b=y=>{be(y,{class:"h-4 w-4 text-green-500"})},S=y=>{$e(y,{class:"h-4 w-4"})};R(g,y=>{e(z)===e(h)?y(b):y(S,!1)})}l(m,w)},$$slots:{default:!0}}),s(p),s(u),I(2),s(_),l(t,d)},$$slots:{default:!0}});var xe=r(fe,2);F(xe,{children:(t,f)=>{var d=ht(),a=r(P(d),2),x=o(a);et(x,{get tabs(){return Pe},variant:"pills",get activeTab(){return e(X)},set activeTab(u){n(X,u,!0)}}),s(a);var _=r(a,2);ge(_,{get code(){return e(ce).code},get language(){return e(ce).language}}),l(t,d)},$$slots:{default:!0}});var Ue=r(xe,2);F(Ue,{children:(t,f)=>{var d=_t(),a=P(d),x=o(a);ct(x,{class:"h-5 w-5 text-primary"}),I(2),s(a);var _=r(a,4),u=r(o(_),2);{let p=k(()=>`import { Outlet } from '@outlet/sdk';

const outlet = new Outlet('YOUR_API_KEY', '${e(h).replace("/sdk/v1","")}');

// Send a transactional email
const result = await outlet.emails.sendEmail({
  to: 'user@example.com',
  subject: 'Welcome!',
  html_body: '<h1>Welcome to our platform</h1>',
});

console.log('Message ID:', result.message_id);`);ge(u,{get code(){return e(p)},language:"typescript"})}s(_),l(t,d)},$$slots:{default:!0}}),s(ee);var je=r(ee,2);Qe(je,{title:"Create API Key",onclose:Z,get show(){return e(N)},set show(f){n(N,f,!0)},footer:f=>{var d=bt(),a=o(d);{var x=u=>{O(u,{type:"primary",onclick:Z,children:(p,c)=>{I();var i=re("Done");l(p,i)},$$slots:{default:!0}})},_=u=>{var p=gt(),c=P(p);O(c,{type:"secondary",onclick:Z,children:(m,$)=>{I();var w=re("Cancel");l(m,w)},$$slots:{default:!0}});var i=r(c,2);{let m=k(()=>e(B)||!e(L).trim());O(i,{type:"primary",onclick:Se,get disabled(){return e(m)},children:($,w)=>{I();var g=re();W(()=>j(g,e(B)?"Creating...":"Create Key")),l($,g)},$$slots:{default:!0}})}l(u,p)};R(a,u=>{e(E)?u(x):u(_,!1)})}s(d),l(f,d)},children:(f,d)=>{var a=q(),x=P(a);{var _=p=>{var c=Pt(),i=o(c);ne(i,{type:"warning",title:"Save your API key",children:(b,S)=>{var y=$t();l(b,y)},$$slots:{default:!0}});var m=r(i,2),$=r(o(m),2),w=o($);le(w,{type:"text",get value(){return e(E)},readonly:!0,class:"font-mono text-sm"});var g=r(w,2);O(g,{type:"secondary",onclick:()=>ie(e(E)),children:(b,S)=>{var y=q(),H=P(y);{var U=C=>{be(C,{class:"h-4 w-4 text-green-500"})},D=C=>{$e(C,{class:"h-4 w-4"})};R(H,C=>{e(z)===e(E)?C(U):C(D,!1)})}l(b,y)},$$slots:{default:!0}}),s($),s(m),s(c),l(p,c)},u=p=>{var c=At(),i=o(c),m=r(o(i),2);le(m,{type:"text",id:"key-name",placeholder:"e.g., Production API, Development",get value(){return e(L)},set value($){n(L,$,!0)}}),I(2),s(i),s(c),l(p,c)};R(x,p=>{e(E)?p(_):p(u,!1)})}l(f,a)},$$slots:{footer:!0,default:!0}}),l(G,pe),Me()}export{Dt as component};
