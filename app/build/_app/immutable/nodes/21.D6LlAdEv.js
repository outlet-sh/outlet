import"../chunks/DsnmJJEf.js";import{a as ot,e as lt,f as O,i as t,t as dt,c as a,h as e,b as H,j as m,s as D,g as it,n as d,r as s}from"../chunks/odI--4Kw.js";import{s as ct}from"../chunks/BNMK9zrr.js";import{a as i,c as nt,t as I,f as y}from"../chunks/CCInJaqd.js";import{i as mt}from"../chunks/COaB7F5h.js";import{h as pt}from"../chunks/hUjpWz_l.js";import{A as ut,C as M,B as xt,e as gt,g as vt,b as ft,f as yt}from"../chunks/Bx02zSD-.js";import{B as A}from"../chunks/BD7k7D4N.js";import{S as ht}from"../chunks/B44UtuBV.js";import{M as bt}from"../chunks/CrfBT87S.js";var $t=y(`<p>Connect your applications to Outlet via SMTP to send transactional or marketing emails. This
			allows you to use Outlet as your SMTP relay server.</p> <p class="mt-2 text-sm"><strong>Workspace scoping:</strong> Your API key determines which brand/organization emails
			are sent from. If you have multiple brands, create separate API keys for each.</p>`,1),_t=y(`<div class="flex items-center gap-3 mb-4"><!> <h2 class="text-lg font-medium text-text">Connection Settings</h2></div> <div class="space-y-4"><div class="grid grid-cols-2 gap-4"><div><label class="form-label">SMTP Host</label> <div class="flex gap-2"><code class="flex-1 bg-surface-tertiary px-3 py-2 rounded text-sm font-mono text-text"> </code> <!></div></div> <div><label class="form-label">SMTP Port</label> <code class="block bg-surface-tertiary px-3 py-2 rounded text-sm font-mono text-text"></code></div></div> <div class="grid grid-cols-2 gap-4"><div><label class="form-label">Username</label> <code class="block bg-surface-tertiary px-3 py-2 rounded text-sm font-mono text-text">api</code> <p class="mt-1 text-xs text-text-muted">Or use your organization slug</p></div> <div><label class="form-label">Password</label> <code class="block bg-surface-tertiary px-3 py-2 rounded text-sm font-mono text-text">YOUR_API_KEY</code> <p class="mt-1 text-xs text-text-muted"><a href="/settings/api" class="text-primary hover:underline">Get an API key</a> from the
						API tab</p></div></div> <div class="pt-4 border-t border-border"><h3 class="font-medium text-text mb-3">Encryption</h3> <div class="flex items-center gap-2"><!> <span class="text-sm text-text-muted">Encryption upgrades automatically when available</span></div></div></div>`,1),Tt=y('<div class="flex items-center gap-3 mb-4"><!> <h2 class="text-lg font-medium text-text">Custom Headers</h2></div> <p class="text-sm text-text-muted mb-4">Use these custom headers to control how Outlet processes your emails.</p> <div class="overflow-x-auto"><table class="w-full text-sm"><thead><tr class="border-b border-border"><th class="text-left py-2 font-medium text-text-muted">Header</th><th class="text-left py-2 font-medium text-text-muted">Values</th><th class="text-left py-2 font-medium text-text-muted">Description</th></tr></thead><tbody class="divide-y divide-border/50"><tr><td class="py-3 font-mono text-primary whitespace-nowrap">X-Outlet-Org</td><td class="py-3"><code class="text-xs bg-surface-tertiary px-1 rounded">slug</code></td><td class="py-3 text-text-muted"><strong>Required.</strong> Organization/brand slug to send emails from. Get this from your brand URL.</td></tr><tr><td class="py-3 font-mono text-primary whitespace-nowrap">X-Outlet-Type</td><td class="py-3"><!> <!></td><td class="py-3 text-text-muted">Email type. Default: <code class="text-xs bg-surface-tertiary px-1 rounded">transactional</code></td></tr><tr><td class="py-3 font-mono text-primary whitespace-nowrap">X-Outlet-List</td><td class="py-3"><code class="text-xs bg-surface-tertiary px-1 rounded">slug</code></td><td class="py-3 text-text-muted">Associate email with a list (for marketing). List slug is scoped to your brand.</td></tr><tr><td class="py-3 font-mono text-primary whitespace-nowrap">X-Outlet-Tags</td><td class="py-3"><code class="text-xs bg-surface-tertiary px-1 rounded">tag1,tag2,tag3</code></td><td class="py-3 text-text-muted">Comma-separated tags to apply to recipient</td></tr><tr><td class="py-3 font-mono text-primary whitespace-nowrap">X-Outlet-Template</td><td class="py-3"><code class="text-xs bg-surface-tertiary px-1 rounded">slug</code></td><td class="py-3 text-text-muted">Use a predefined email template. Template slug is scoped to your brand.</td></tr><tr><td class="py-3 font-mono text-primary whitespace-nowrap">X-Outlet-Track</td><td class="py-3"><code class="text-xs bg-surface-tertiary px-1 rounded">opens,clicks</code> <code class="text-xs bg-surface-tertiary px-1 rounded">none</code></td><td class="py-3 text-text-muted">Enable/disable tracking. Default: <code class="text-xs bg-surface-tertiary px-1 rounded">opens,clicks</code></td></tr><tr><td class="py-3 font-mono text-primary whitespace-nowrap">X-Outlet-Meta-*</td><td class="py-3"><code class="text-xs bg-surface-tertiary px-1 rounded">any value</code></td><td class="py-3 text-text-muted">Custom metadata. E.g., <code class="text-xs bg-surface-tertiary px-1 rounded">X-Outlet-Meta-Order-ID: 12345</code></td></tr></tbody></table></div>',1),Pt=y('<h2 class="text-lg font-medium text-text mb-4">Code Examples</h2> <div class="mb-4"><!></div> <!>',1),Ot=y(`<h2 class="text-lg font-medium text-text mb-2">SMTP Limits</h2> <p class="text-sm text-text-muted mb-4">These are the default limits for emails sent via the SMTP ingress server. Emails exceeding
			these limits will be rejected.</p> <div class="grid grid-cols-2 gap-4 text-sm"><div class="flex justify-between p-3 bg-surface-secondary rounded"><span class="text-text-muted">Max message size</span> <span class="text-text font-medium">25 MB</span></div> <div class="flex justify-between p-3 bg-surface-secondary rounded"><span class="text-text-muted">Max recipients per message</span> <span class="text-text font-medium">100</span></div></div>`,1),Mt=y('<div class="space-y-6"><!> <!> <!> <!> <!></div>');function Ut(F){let z=ot(window.location.hostname),w=D("");function K(r){navigator.clipboard.writeText(r),H(w,r,!0),setTimeout(()=>H(w,""),2e3)}const c=m(()=>z),h=587,W=[{id:"curl",label:"cURL"},{id:"nodejs",label:"Node.js"},{id:"python",label:"Python"},{id:"go",label:"Go"},{id:"php",label:"PHP"}];let S=D("curl");const q=m(()=>`# Send email via SMTP using curl
curl --url "smtp://${t(c)}:${h}" \\
  --ssl-reqd \\
  --user "api:YOUR_API_KEY" \\
  --mail-from "you@example.com" \\
  --mail-rcpt "recipient@example.com" \\
  --upload-file - << EOF
From: you@example.com
To: recipient@example.com
Subject: Hello from Outlet
Content-Type: text/html
X-Outlet-Org: your-org-slug
X-Outlet-Type: transactional
X-Outlet-Track: opens,clicks
X-Outlet-Meta-User-ID: 12345

<h1>Welcome!</h1>
<p>This is a test email sent via SMTP.</p>
EOF`),G=m(()=>`import nodemailer from 'nodemailer';

const transporter = nodemailer.createTransport({
  host: '${t(c)}',
  port: ${h},
  secure: false, // STARTTLS
  auth: {
    user: 'api',
    pass: 'YOUR_API_KEY'
  }
});

await transporter.sendMail({
  from: 'you@example.com',
  to: 'recipient@example.com',
  subject: 'Hello from Outlet',
  html: '<h1>Welcome!</h1>',
  headers: {
    'X-Outlet-Org': 'your-org-slug',
    'X-Outlet-Type': 'transactional',
    'X-Outlet-Track': 'opens,clicks',
    'X-Outlet-Meta-User-ID': '12345'
  }
});`),N=m(()=>`import smtplib
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart

msg = MIMEMultipart('alternative')
msg['From'] = 'you@example.com'
msg['To'] = 'recipient@example.com'
msg['Subject'] = 'Hello from Outlet'
msg['X-Outlet-Org'] = 'your-org-slug'
msg['X-Outlet-Type'] = 'transactional'
msg['X-Outlet-Track'] = 'opens,clicks'

msg.attach(MIMEText('<h1>Welcome!</h1>', 'html'))

with smtplib.SMTP('${t(c)}', ${h}) as server:
    server.starttls()
    server.login('api', 'YOUR_API_KEY')
    server.sendmail(msg['From'], msg['To'], msg.as_string())`),V=m(()=>`package main

import (
    "net/smtp"
)

func main() {
    auth := smtp.PlainAuth("", "api", "YOUR_API_KEY", "${t(c)}")

    msg := []byte("From: you@example.com\\r\\n" +
        "To: recipient@example.com\\r\\n" +
        "Subject: Hello from Outlet\\r\\n" +
        "X-Outlet-Org: your-org-slug\\r\\n" +
        "X-Outlet-Type: transactional\\r\\n" +
        "Content-Type: text/html\\r\\n" +
        "\\r\\n" +
        "<h1>Welcome!</h1>")

    err := smtp.SendMail("${t(c)}:${h}", auth,
        "you@example.com", []string{"recipient@example.com"}, msg)
}`),J=m(()=>`use PHPMailer\\PHPMailer\\PHPMailer;

$mail = new PHPMailer(true);

$mail->isSMTP();
$mail->Host = '${t(c)}';
$mail->Port = ${h};
$mail->SMTPAuth = true;
$mail->Username = 'api';
$mail->Password = 'YOUR_API_KEY';
$mail->SMTPSecure = PHPMailer::ENCRYPTION_STARTTLS;

$mail->setFrom('you@example.com');
$mail->addAddress('recipient@example.com');
$mail->Subject = 'Hello from Outlet';
$mail->isHTML(true);
$mail->Body = '<h1>Welcome!</h1>';

$mail->addCustomHeader('X-Outlet-Org', 'your-org-slug');
$mail->addCustomHeader('X-Outlet-Type', 'transactional');
$mail->addCustomHeader('X-Outlet-Track', 'opens,clicks');

$mail->send();`);let Q=m(()=>({curl:{code:t(q),language:"bash"},nodejs:{code:t(G),language:"javascript"},python:{code:t(N),language:"python"},go:{code:t(V),language:"go"},php:{code:t(J),language:"php"}})),j=m(()=>t(Q)[t(S)]);var k=Mt();pt("1suqzf",r=>{lt(()=>{it.title="SMTP - Settings"})});var U=a(k);ut(U,{type:"info",title:"SMTP Ingress",children:(r,_)=>{var o=$t();d(2),i(r,o)},$$slots:{default:!0}});var Y=e(U,2);M(Y,{children:(r,_)=>{var o=_t(),l=O(o),b=a(l);ht(b,{class:"h-5 w-5 text-primary"}),d(2),s(l);var p=e(l,2),n=a(p),u=a(n),$=e(a(u),2),x=a($),T=a(x,!0);s(x);var C=e(x,2);xt(C,{type:"secondary",size:"sm",onclick:()=>K(t(c)),children:(E,et)=>{var P=nt(),at=O(P);{var st=f=>{ft(f,{class:"h-4 w-4 text-green-500"})},rt=f=>{yt(f,{class:"h-4 w-4"})};mt(at,f=>{t(w)===t(c)?f(st):f(rt,!1)})}i(E,P)},$$slots:{default:!0}}),s($),s(u);var g=e(u,2),X=e(a(g),2);X.textContent="587",s(g),s(n);var v=e(n,4),B=e(a(v),2),tt=a(B);A(tt,{type:"success",children:(E,et)=>{d();var P=I("STARTTLS");i(E,P)},$$slots:{default:!0}}),d(2),s(B),s(v),s(p),dt(()=>ct(T,t(c))),i(r,o)},$$slots:{default:!0}});var R=e(Y,2);M(R,{children:(r,_)=>{var o=Tt(),l=O(o),b=a(l);bt(b,{class:"h-5 w-5 text-primary"}),d(2),s(l);var p=e(l,4),n=a(p),u=e(a(n)),$=e(a(u)),x=e(a($)),T=a(x);A(T,{type:"secondary",children:(g,X)=>{d();var v=I("transactional");i(g,v)},$$slots:{default:!0}});var C=e(T,2);A(C,{type:"secondary",children:(g,X)=>{d();var v=I("marketing");i(g,v)},$$slots:{default:!0}}),s(x),d(),s($),d(5),s(u),s(n),s(p),i(r,o)},$$slots:{default:!0}});var L=e(R,2);M(L,{children:(r,_)=>{var o=Pt(),l=e(O(o),2),b=a(l);gt(b,{get tabs(){return W},variant:"pills",get activeTab(){return t(S)},set activeTab(n){H(S,n,!0)}}),s(l);var p=e(l,2);vt(p,{get code(){return t(j).code},get language(){return t(j).language}}),i(r,o)},$$slots:{default:!0}});var Z=e(L,2);M(Z,{children:(r,_)=>{var o=Ot();d(4),i(r,o)},$$slots:{default:!0}}),s(k),i(F,k)}export{Ut as component};
