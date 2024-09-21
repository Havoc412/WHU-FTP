import { defineStore } from "pinia";

export default defineStore("routeStore", {
    state: () : {cacheRouteList: string[] } => {
        return {
          // 默认需要缓存的界面
          cacheRouteList: ["Login"], //'IllustrateView' 需要 动态添加
        };
    },
    actions: {
        addCacheRoute(name: string) {
            this.cacheRouteList.push(name);
        },
        removeCacheRoute(name: string) {
            for(let i = this.cacheRouteList.length - 1; i>=0; i--)
                if(this.cacheRouteList[i] === name)
                    this.cacheRouteList.splice(i, 1);
        }
    }
});