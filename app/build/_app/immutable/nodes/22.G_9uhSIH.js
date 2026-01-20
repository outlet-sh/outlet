import"../chunks/DsnmJJEf.js";import{a as ot,e as lt,h as e,c as a,f as O,i as t,j as n,t as it,b as H,s as D,g as dt,n as i,r as s}from"../chunks/DDjrJUyi.js";import{s as ct}from"../chunks/CwHIudaz.js";import{f as y,a as d,c as mt,t as I}from"../chunks/DbkN8AN9.js";import{i as nt}from"../chunks/CGvWOYBA.js";import{h as pt}from"../chunks/DHP91Nzo.js";import{A as ut,C as w,B as xt,f as gt,h as vt,b as ft,g as yt}from"../chunks/DkmOYVX9.js";import{B as A}from"../chunks/42UF-HRp.js";import{S as ht}from"../chunks/DQg3mEaZ.js";import{M as bt}from"../chunks/CxHuRNbG.js";var $t=y(`<p>Connect your applications to Outlet via SMTP to send transactional or marketing emails. This
			allows you to use Outlet as your SMTP relay server.</p> <p class="mt-2 text-sm"><strong>Workspace scoping:</strong> Your API key determines which workspace/organization emails
			are sent from. If you have multiple workspaces, create separate API keys for each.</p>`,1),_t=y(`<div class="flex items-center gap-3 mb-4"><!> <h2 class="text-lg font-medium text-text">Connection Settings</h2></div> <div class="space-y-4"><div class="grid grid-cols-2 gap-4"><div><label class="form-label">SMTP Host</label> <div class="flex gap-2"><code class="flex-1 bg-surface-tertiary px-3 py-2 rounded text-sm font-mono text-text"> </code> <!></div></div> <div><label class="form-label">SMTP Port</label> <code class="block bg-surface-tertiary px-3 py-2 rounded text-sm font-mono text-text"></code></div></div> <div class="grid grid-cols-2 gap-4"><div><label class="form-label">Username</label> <code class="block bg-surface-tertiary px-3 py-2 rounded text-sm font-mono text-text">api</code> <p class="mt-1 text-xs text-text-muted">Or use your organization slug</p></div> <div><label class="form-label">Password</label> <code class="block bg-surface-tertiary px-3 py-2 rounded text-sm font-mono text-text">YOUR_API_KEY</code> <p class="mt-1 text-xs text-text-muted"><a href="/settings/api" class="text-primary hover:underline">Get an API key</a> from the
						API tab</p></div></div> <div class="pt-4 border-t border-border"><h3 class="font-medium text-text mb-3">Encryption</h3> <div class="flex items-center gap-2"><!> <span class="text-sm text-text-muted">Encryption upgrades automatically when available</span></div></div></div>`,1),Tt=y('<div class="flex items-center gap-3 mb-4"><!> <h2 class="text-lg font-medium text-text">Custom Headers</h2></div> <p class="text-sm text-text-muted mb-4">Use these custom headers to control how Outlet processes your emails.</p> <div class="overflow-x-auto"><table class="w-full text-sm"><thead><tr class="border-b border-border"><th class="text-left py-2 font-medium text-text-muted">Header</th><th class="text-left py-2 font-medium text-text-muted">Values</th><th class="text-left py-2 font-medium text-text-muted">Description</th></tr></thead><tbody class="divide-y divide-border/50"><tr><td class="py-3 font-mono text-primary whitespace-nowrap">X-Outlet-Org</td><td class="py-3"><code class="text-xs bg-surface-tertiary px-1 rounded">slug</code></td><td class="py-3 text-text-muted"><strong>Required.</strong> Organization/workspace slug to send emails from. Get this from your workspace URL.</td></tr><tr><td class="py-3 font-mono text-primary whitespace-nowrap">X-Outlet-Type</td><td class="py-3"><!> <!></td><td class="py-3 text-text-muted">Email type. Default: <code class="text-xs bg-surface-tertiary px-1 rounded">transactional</code></td></tr><tr><td class="py-3 font-mono text-primary whitespace-nowrap">X-Outlet-List</td><td class="py-3"><code class="text-xs bg-surface-tertiary px-1 rounded">slug</code></td><td class="py-3 text-text-muted">Associate email with a list (for marketing). List slug is scoped to your workspace.</td></tr><tr><td class="py-3 font-mono text-primary whitespace-nowrap">X-Outlet-Tags</td><td class="py-3"><code class="text-xs bg-surface-tertiary px-1 rounded">tag1,tag2,tag3</code></td><td class="py-3 text-text-muted">Comma-separated tags to apply to recipient</td></tr><tr><td class="py-3 font-mono text-primary whitespace-nowrap">X-Outlet-Template</td><td class="py-3"><code class="text-xs bg-surface-tertiary px-1 rounded">slug</code></td><td class="py-3 text-text-muted">Use a predefined email template. Template slug is scoped to your workspace.</td></tr><tr><td class="py-3 font-mono text-primary whitespace-nowrap">X-Outlet-Track</td><td class="py-3"><code class="text-xs bg-surface-tertiary px-1 rounded">opens,clicks</code> <code class="text-xs bg-surface-tertiary px-1 rounded">none</code></td><td class="py-3 text-text-muted">Enable/disable tracking. Default: <code class="text-xs bg-surface-tertiary px-1 rounded">opens,clicks</code></td></tr><tr><td class="py-3 font-mono text-primary whitespace-nowrap">X-Outlet-Meta-*</td><td class="py-3"><code class="text-xs bg-surface-tertiary px-1 rounded">any value</code></td><td class="py-3 text-text-muted">Custom metadata. E.g., <code class="text-xs bg-surface-tertiary px-1 rounded">X-Outlet-Meta-Order-ID: 12345</code></td></tr></tbody></table></div>',1),Pt=y('<h2 class="text-lg font-medium text-text mb-4">Code Examples</h2> <div class="mb-4"><!></div> <!>',1),Ot=y(`<h2 class="text-lg font-medium text-text mb-2">SMTP Limits</h2> <p class="text-sm text-text-muted mb-4">These are the default limits for emails sent via the SMTP ingress server. Emails exceeding
			these limits will be rejected.</p> <div class="grid grid-cols-2 gap-4 text-sm"><div class="flex justify-between p-3 bg-surface-secondary rounded"><span class="text-text-muted">Max message size</span> <span class="text-text font-medium">25 MB</span></div> <div class="flex justify-between p-3 bg-surface-secondary rounded"><span class="text-text-muted">Max recipients per message</span> <span class="text-text font-medium">100</span></div></div>`,1),wt=y('<div class="space-y-6"><!> <!> <!> <!> <!></div>');function Ut(F){let K=ot(window.location.hostname),M=D("");function W(r){navigator.clipboard.writeText(r),H(M,r,!0),setTimeout(()=>H(M,""),2e3)}const c=n(()=>K),h=587,z=[{id:"curl",label:"cURL"},{id:"nodejs",label:"Node.js"},{id:"python",label:"Python"},{id:"go",label:"Go"},{id:"php",label:"PHP"}];let k=D("curl");const G=n(()=>`# Send email via SMTP using curl
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
EOF`),N=n(()=>`import nodemailer from 'nodemailer';

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
});`),q=n(()=>`import smtplib
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
    server.sendmail(msg['From'], msg['To'], msg.as_string())`),V=n(()=>`package main

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
}`),J=n(()=>`use PHPMailer\\PHPMailer\\PHPMailer;

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

$mail->send();`);let Q=n(()=>({curl:{code:t(G),language:"bash"},nodejs:{code:t(N),language:"javascript"},python:{code:t(q),language:"python"},go:{code:t(V),language:"go"},php:{code:t(J),language:"php"}})),j=n(()=>t(Q)[t(k)]);var S=wt();pt("ekcsaa",r=>{lt(()=>{dt.title="SMTP - Settings"})});var U=a(S);ut(U,{type:"info",title:"SMTP Ingress",children:(r,_)=>{var o=$t();i(2),d(r,o)},$$slots:{default:!0}});var Y=e(U,2);w(Y,{children:(r,_)=>{var o=_t(),l=O(o),b=a(l);ht(b,{class:"h-5 w-5 text-primary"}),i(2),s(l);var p=e(l,2),m=a(p),u=a(m),$=e(a(u),2),x=a($),T=a(x,!0);s(x);var C=e(x,2);xt(C,{type:"secondary",size:"sm",onclick:()=>W(t(c)),children:(E,et)=>{var P=mt(),at=O(P);{var st=f=>{ft(f,{class:"h-4 w-4 text-green-500"})},rt=f=>{yt(f,{class:"h-4 w-4"})};nt(at,f=>{t(M)===t(c)?f(st):f(rt,!1)})}d(E,P)},$$slots:{default:!0}}),s($),s(u);var g=e(u,2),X=e(a(g),2);X.textContent="587",s(g),s(m);var v=e(m,4),B=e(a(v),2),tt=a(B);A(tt,{type:"success",children:(E,et)=>{i();var P=I("STARTTLS");d(E,P)},$$slots:{default:!0}}),i(2),s(B),s(v),s(p),it(()=>ct(T,t(c))),d(r,o)},$$slots:{default:!0}});var R=e(Y,2);w(R,{children:(r,_)=>{var o=Tt(),l=O(o),b=a(l);bt(b,{class:"h-5 w-5 text-primary"}),i(2),s(l);var p=e(l,4),m=a(p),u=e(a(m)),$=e(a(u)),x=e(a($)),T=a(x);A(T,{type:"secondary",children:(g,X)=>{i();var v=I("transactional");d(g,v)},$$slots:{default:!0}});var C=e(T,2);A(C,{type:"secondary",children:(g,X)=>{i();var v=I("marketing");d(g,v)},$$slots:{default:!0}}),s(x),i(),s($),i(5),s(u),s(m),s(p),d(r,o)},$$slots:{default:!0}});var L=e(R,2);w(L,{children:(r,_)=>{var o=Pt(),l=e(O(o),2),b=a(l);gt(b,{get tabs(){return z},variant:"pills",get activeTab(){return t(k)},set activeTab(m){H(k,m,!0)}}),s(l);var p=e(l,2);vt(p,{get code(){return t(j).code},get language(){return t(j).language}}),d(r,o)},$$slots:{default:!0}});var Z=e(L,2);w(Z,{children:(r,_)=>{var o=Ot();i(4),d(r,o)},$$slots:{default:!0}}),s(S),d(F,S)}export{Ut as component};
