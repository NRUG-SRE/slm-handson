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

  // New Relic RUM設定スクリプトを環境変数から動的生成
  const newRelicScript = newRelicBrowserKey && newRelicAccountId && newRelicApplicationId ? `
    window.NREUM||(NREUM={});
    NREUM.init={
      distributed_tracing:{enabled:true},
      privacy:{cookies_enabled:true},
      ajax:{deny_list:["bam.nr-data.net"]}
    };

    NREUM.loader_config={
      accountID:"${newRelicAccountId}",
      trustKey:"${newRelicAccountId}",
      agentID:"${newRelicApplicationId}",
      licenseKey:"${newRelicBrowserKey}",
      applicationID:"${newRelicApplicationId}"
    };

    NREUM.info={
      beacon:"bam.nr-data.net",
      errorBeacon:"bam.nr-data.net",
      licenseKey:"${newRelicBrowserKey}",
      applicationID:"${newRelicApplicationId}",
      sa:1
    };

    // New Relicローダースクリプトの動的読み込み
    (function() {
      var script = document.createElement('script');
      script.type = 'text/javascript';
      script.async = true;
      script.src = 'https://js-agent.newrelic.com/nr-loader-spa-1.293.0.min.js';
      
      var firstScript = document.getElementsByTagName('script')[0];
      firstScript.parentNode.insertBefore(script, firstScript);
    })();
  ` : ''

  return (
    <html lang="ja">
      <head>
        {newRelicScript && (
          <Script
            id="new-relic-rum"
            strategy="beforeInteractive"
            dangerouslySetInnerHTML={{
              __html: newRelicScript
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