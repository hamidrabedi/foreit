import React from "react";
import { Refine } from "@refinedev/core";
import { RefineKbar, RefineKbarProvider } from "@refinedev/kbar";
import routerProvider from "@refinedev/react-router-v6";
import { BrowserRouter, Route, Routes, Outlet } from "react-router-dom";
import {
  ThemedLayoutV2,
  ErrorComponent,
  notificationProvider,
  RefineThemes,
} from "@refinedev/antd";
import { App as AntdApp, ConfigProvider } from "antd";
import "@refinedev/antd/dist/reset.css";

import { gogoDataProvider } from "./dataProvider";

// Example pages (you'll need to create these)
// import { UserList, UserCreate, UserEdit, UserShow } from "./pages/users";
// import { ArticleList, ArticleCreate, ArticleEdit, ArticleShow } from "./pages/articles";

function App() {
  const API_URL = "http://localhost:8080/admin/api";

  return (
    <BrowserRouter>
      <RefineKbarProvider>
        <ConfigProvider theme={RefineThemes.Blue}>
          <AntdApp>
            <Refine
              routerProvider={routerProvider}
              dataProvider={gogoDataProvider({
                apiUrl: API_URL,
                getAccessToken: () => {
                  // Get token from your auth system
                  return localStorage.getItem("token");
                },
              })}
              resources={[
                {
                  name: "users",
                  list: "/users",
                  create: "/users/create",
                  edit: "/users/edit/:id",
                  show: "/users/show/:id",
                  meta: {
                    canDelete: true,
                  },
                },
                // Add more resources based on your Gogo models
                // {
                //   name: "articles",
                //   list: "/articles",
                //   create: "/articles/create",
                //   edit: "/articles/edit/:id",
                //   show: "/articles/show/:id",
                // },
              ]}
              notificationProvider={notificationProvider}
              options={{
                syncWithLocation: true,
                warnWhenUnsavedChanges: true,
                useNewQueryKeys: true,
                projectId: "gogo-admin",
              }}
            >
              <Routes>
                <Route
                  element={
                    <ThemedLayoutV2>
                      <Outlet />
                    </ThemedLayoutV2>
                  }
                >
                  <Route
                    index
                    element={<div>Welcome to Gogo Admin</div>}
                  />
                  {/* Example routes - uncomment when you create the pages */}
                  {/* <Route path="/users">
                    <Route index element={<UserList />} />
                    <Route path="create" element={<UserCreate />} />
                    <Route path="edit/:id" element={<UserEdit />} />
                    <Route path="show/:id" element={<UserShow />} />
                  </Route> */}
                  <Route path="*" element={<ErrorComponent />} />
                </Route>
              </Routes>
              <RefineKbar />
            </Refine>
          </AntdApp>
        </ConfigProvider>
      </RefineKbarProvider>
    </BrowserRouter>
  );
}

export default App;

