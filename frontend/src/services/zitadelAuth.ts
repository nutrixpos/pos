import { createZITADELAuth } from "zitadel-vue/lib/src";
import { User } from "oidc-client";
import { SignInType, type OidcAuth } from 'vue-oidc-client/vue3';

const zitadelAuth: { oidcAuth: OidcAuth; hasRole: (role: string) => any } = createZITADELAuth({
    project_resource_id: process.env.VUE_APP_ZITADEL_PROJECT_RESOURCE_ID,
    client_id: process.env.VUE_APP_ZITADEL_CLIENT_ID,
    issuer: process.env.VUE_APP_ZITADEL_ISSUER,
},"Zitadel",SignInType.Window,"http://localhost:8080/",{
   post_logout_redirect_uri:'http://localhost:8080/',
})

// handle events
zitadelAuth.oidcAuth.events.addAccessTokenExpiring(function() {
  // eslint-disable-next-line no-console
  console.log('access token expiring')
})

zitadelAuth.oidcAuth.events.addAccessTokenExpired(function() {
  // eslint-disable-next-line no-console
  console.log('access token expired')
})

zitadelAuth.oidcAuth.events.addSilentRenewError(function(err: Error) {
  // eslint-disable-next-line no-console
  console.error('silent renew error', err)
})

zitadelAuth.oidcAuth.events.addUserLoaded(function(user: User) {
  // eslint-disable-next-line no-console
  console.log('user loaded', user)
})

zitadelAuth.oidcAuth.events.addUserUnloaded(function() {
  // eslint-disable-next-line no-console
  console.log('user unloaded')
})

zitadelAuth.oidcAuth.events.addUserSignedOut(function() {
  // eslint-disable-next-line no-console
  console.log('user signed out')
})

zitadelAuth.oidcAuth.events.addUserSessionChanged(function() {
  // eslint-disable-next-line no-console
  console.log('user session changed')
})

export default zitadelAuth