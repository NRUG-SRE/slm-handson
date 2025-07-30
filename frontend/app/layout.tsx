import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'
import Header from '@/components/layout/Header'
import Footer from '@/components/layout/Footer'
import Script from 'next/script'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'SLM ハンズオン ECサイト',
  description: 'New Relic Service Level Management ハンズオン用のECサイトデモ',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const newRelicBrowserKey = process.env.NEXT_PUBLIC_NEW_RELIC_BROWSER_KEY
  const newRelicAccountId = process.env.NEXT_PUBLIC_NEW_RELIC_ACCOUNT_ID
  const newRelicApplicationId = process.env.NEXT_PUBLIC_NEW_RELIC_APPLICATION_ID

  return (
    <html lang="ja">
      <head>
        {newRelicBrowserKey && newRelicAccountId && newRelicApplicationId && (
          <Script
            strategy="beforeInteractive"
            dangerouslySetInnerHTML={{
              __html: `
                window.NREUM||(NREUM={});NREUM.info={
                  "beacon":"bam.nr-data.net",
                  "licenseKey":"${newRelicBrowserKey}",
                  "applicationID":"${newRelicApplicationId}",
                  "transactionName":"",
                  "queueTime":0,
                  "applicationTime":0,
                  "agent":"js-agent.newrelic.com/nr-1.260.1.js"
                };
                (function(a,b,c,d,e,f,g){a.NREUM||(a.NREUM={}),
                b in a.NREUM||(a.NREUM[b]=[]);var h=typeof d==="function";
                a.NREUM[b].push({name:c,params:h?null:d,callback:h?d:null,end:Date.now()});
                if(a.readyState==="complete"||a.readyState==="interactive")
                f();else a.addEventListener("DOMContentLoaded",f,false)})
                (document,"addEventListener","DOMContentLoaded",function(){
                (function(a,b){var c=a.createElement(b),d=a.getElementsByTagName(b)[0];
                c.src="https://js-agent.newrelic.com/nr-1.260.1.js";
                c.setAttribute("data-account-id","${newRelicAccountId}");
                c.setAttribute("data-application-id","${newRelicApplicationId}");
                d.parentNode.insertBefore(c,d)})(document,"script")});
              `,
            }}
          />
        )}
      </head>
      <body className={inter.className}>
        <div className="min-h-screen flex flex-col">
          <Header />
          <main className="flex-grow">
            {children}
          </main>
          <Footer />
        </div>
      </body>
    </html>
  )
}