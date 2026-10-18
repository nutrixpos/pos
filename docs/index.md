---
# https://vitepress.dev/reference/default-theme-home-page
layout: home

hero:
  name: "NutrixPOS"
  text: "Free Open Source Restaurant Management System"
  tagline:
  image:
    src: "https://elmawardy.sirv.com/nutrixpos-docs/nutrixdocs0.png"
    alt: nutrix-icon
  actions:
    - theme: brand
      text: Download now - it's free ↗
      link: https://github.com/nutrixpos/pos/releases/latest
    - theme: alt
      text: Meet an advisor
      link: /support

features:
  - title: Cashier
    icon: 🖨
    details: Seamless order creation with high level of control to inventory consumption
  - title: Kitchen
    icon: 🍳
    details: On demand simple interface for displaying and finishing orders, to be used inside the kitchen
  - title: Inventory
    icon: 🗄️
    details: Easily manage your inventory and stay notified about low stock
  - title: Reports
    icon: 📈
    details: Observe inventory consumptions, sales and profit per day down to the product's components level.
---


<div class="grid pt-4 mt-4 px-0">
  <div class="col-12 lg:col-5 lg:pr-6 lg:mt-8 lg:pt-4">
    <h1><span class="gradiented">Intuitive</span> order creation experience</h1>
    <h4 class="mt-3" style="color:var(--vp-c-text-2)">
      Specifically built for restaurants, which eliminates the unwanted complexities and helps you focus on what really matters.
    </h4>
    <ol>
      <li>
        <h5 style="color: var(--vp-c-text-2);">
          Easily divide your products to categories
        </h5>
      </li>
      <li>
        <h5 style="color: var(--vp-c-text-2);">
          Know the potential of your inventory
        </h5>
      </li>
      <li>
        <h5 style="color: var(--vp-c-text-2);">
          Draft current order and continue with other customers
        </h5>
      </li>
      <li>
        <h5 style="color: var(--vp-c-text-2);">
          Receive notifications on low stock
        </h5>
      </li>
    </ol>
  </div>
  <div class="col-12 lg:col-7 mt-3 lg:mt-0 flex justify-content-center align-items-center">
    <Image src="https://elmawardy.sirv.com/nutrixpos-docs/nutrixdocs1.png" alt="Image" preview />
  </div>
</div>

<div class="grid pt-0 pb-4 px-0">
  <div class="col-12 lg:col-7 mt-3 lg:mt-0 flex justify-content-center align-items-center">
    <Image src="https://elmawardy.sirv.com/nutrixpos-docs/Nutrix3.png" alt="Image" preview />
  </div>
  <div class="col-12 lg:col-5 lg:pl-6 lg:pt-3 lg:mt-8 lg:pt-4">
    <h1>Administrate with <span class="gradiented">Ease</span></h1>
    <h4 class="mt-3" style="color:var(--vp-c-text-2)">
      Easily manage your restaurant inventory and get insights about your sales
    </h4>
    <ul>
      <li>
        <h5 style="color: var(--vp-c-text-2);">
          Observe your cost/sales trend
        </h5>
      </li>
      <li>
        <h5 style="color: var(--vp-c-text-2);">
          Know your best selling products
        </h5>
      </li>
      <li>
        <h5 style="color: var(--vp-c-text-2);">
          Manage your inventory and set up alerts on low stock
        </h5>
      </li>
      <li>
        <h5 style="color: var(--vp-c-text-2);">
          Get exact sales per day and know how much profit you made
        </h5>
      </li>
    </ul>
  </div>
</div>

<Card class="mb-6 mt-8 nutrixhub-section">
    <template #content>
        <div class="m-0 grid">
            <div class="lg:col-6 col-12 flex flex-column justify-content-center" style="color:#334155">
              <h3 style="font-size:2rem;color:#334155;">NutrixHub</h3>
              <p class="p-1 mt-3" style="color:#334155">
                Powerful SaaS platform that turns data from your NutrixPOS instances into actionable growth insights. Get AI-driven sales offers, monitor real-time performance and inventory, and access dedicated support—all from one centralized hub.
              </p>
              <a target="_blank" style="color:#334155" class="flex align-items-center" href="https://hub.nutrixpos.com">Know more &nbsp;<span class="pi pi-external-link"></span></a>
            </div>
            <div class="lg:col-6 col-12">
              <Image src="https://elmawardy.sirv.com/nutrixhub/nutrixhub.png" alt="Image" preview />
            </div>
        </div>
    </template>
</Card>



<div class="grid md:p-5 mt-8">
  <div class="col-12 md:mt-4 md:pt-4 md:col-6">
      <h1>Stay<span class="gradiented"> free</span> stay open source</h1>
      <h5 class="mt-3" style="color:var(--vp-c-text-2)">
          Open source is more than free software—it’s about building a community that works together to create something meaningful for everyone.
      </h5>
  </div>
  <div class="col-12 flex m-0 mt-2 p-0">
    <a href="https://github.com/nutrixpos/">
      <Button style="background-color: var(--vp-c-brand); border-color:var(--vp-c-brand);color:var(--vp-button-brand-text)" rounded label="Contribute" icon="pi pi-github" size="small" class="mx-1"/>
    </a>
    <a href="https://t.me/nutrixpos">
      <Button style="color:#3390EC;background-color:white;" severity="secondary" rounded label="Connect" icon="pi pi-telegram" size="small" class="mx-1"/>
    </a>
    <a href="/sponsor.html">
      <Button style="color:#F15669;background-color:white;" severity="secondary" rounded label="Sponsor" icon="pi pi-heart" size="small" class="mx-1"/>
    </a>
  </div>
</div>

