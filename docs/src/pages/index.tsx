import React from 'react';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import AsciinemaPlayer from "@site/src/components/AsciinemaPlayer";


export default function Home(): JSX.Element {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title={`Hello from ${siteConfig.title}`}
      description={siteConfig.tagline}>
      <main style={{width: "80%", margin:"0 auto"}}>
          <h1>Seki</h1>
          <AsciinemaPlayer src="terminal.cast" preload={true} loop={true} autoPlay={true} terminalFontSize={"16px"} fit={false}/>
      </main>
    </Layout>
  );
}
