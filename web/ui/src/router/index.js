import Vue from 'vue'
import Router from 'vue-router'
import store from '@/store'
import Index from '@/components/Index'
import Login from '@/components/Login'
import Dashboard from '@/components/Dashboard'
import Servers from '@/components/Servers'
import Users from '@/components/Users'
import Sessions from '@/components/Sessions'
import Replay from '@/components/Replay'
import Settings from '@/components/Settings'

// children of settings
import Profile from '@/components/settings/Profile'
import ChangePassword from '@/components/settings/ChangePassword'
import Keys from '@/components/settings/Keys'
import Tokens from '@/components/settings/Tokens'

// children of users
import UserDetail from '@/components/UserDetail'

Vue.use(Router)

let router = new Router({
  routes: [
    {
      path: '/',
      name: 'Index',
      component: Index
    },
    {
      path: '/login',
      name: 'Login',
      component: Login
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: Dashboard,
      meta: {
        requiresLoggedIn: true
      }
    },
    {
      path: '/servers',
      name: 'Servers',
      component: Servers,
      meta: {
        requiresLoggedInAsAdmin: true
      }
    },
    {
      path: '/users',
      name: 'Users',
      component: Users,
      meta: {
        requiresLoggedInAsAdmin: true
      }
    },
    {
      path: '/users/:account',
      name: 'UserDetail',
      component: UserDetail,
      meta: {
        requiresLoggedInAsAdmin: true
      }
    },
    {
      path: '/sessions',
      name: 'Sessions',
      component: Sessions,
      meta: {
        requiresLoggedInAsAdmin: true
      }
    },
    {
      path: '/sessions/:id/replay',
      name: 'Replay',
      component: Replay,
      meta: {
        requiresLoggedInAsAdmin: true,
        hidesNavigationBar: true
      }
    },
    {
      path: '/settings',
      name: 'Settings',
      component: Settings,
      meta: {
        requiresLoggedInAsAdmin: true
      },
      children: [
        {
          path: 'profile',
          name: 'Profile',
          component: Profile
        },
        {
          path: 'change-password',
          name: 'ChangePassword',
          component: ChangePassword
        },
        {
          path: 'keys',
          name: 'Keys',
          component: Keys
        },
        {
          path: 'tokens',
          name: 'Tokens',
          component: Tokens
        }
      ]
    }
  ]
})

router.beforeEach((to, from, next) => {
  if (
    (to.meta.requiresLoggedIn || to.meta.requiresLoggedInAsAdmin) &&
    !store.getters.isLoggedIn
  ) {
    next('/login')
  } else if (
    to.meta.requiresLoggedInAsAdmin &&
    !store.getters.isLoggedInAsAdmin
  ) {
    next('/dashboard')
  } else {
    next(true)
  }
})

export default router
