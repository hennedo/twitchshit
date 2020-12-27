import Vue from 'vue'
import VueRouter from 'vue-router'
import Overlay from '../view/Overlay'
import TwitchOauth from "../view/TwitchOauth";
import OverlayPause from "../view/OverlayPause";
import OverlayFullscreen from "../view/OverlayFullscreen";

Vue.use(VueRouter)

const routes = [
    {
        path: '/',
        component: Overlay
    },
    {
        path: '/fullscreen',
        component: OverlayFullscreen
    },
    {
        path: '/pause',
        component: OverlayPause
    },
    {
        path: '/oauth/twitch',
        component: TwitchOauth,
        name: 'twitchOauth'
    }
]

const router = new VueRouter({
    mode: 'history',
    base: '/',
    activeClass: 'active',
    routes
})

export default router