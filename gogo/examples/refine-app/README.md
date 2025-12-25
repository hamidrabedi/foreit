# Refine + Gogo Example App

This is an example React application using Refine with the Gogo admin engine backend.

## Setup

1. Install dependencies:

```bash
npm install
# or
yarn install
```

2. Start your Gogo backend server (see main README)

3. Update the API URL in `src/App.tsx`:

```typescript
const API_URL = "http://localhost:8080/admin/api";
```

4. Start the development server:

```bash
npm start
# or
yarn start
```

## Features

- Full CRUD operations for registered models
- Authentication (if configured)
- Pagination
- Sorting
- Filtering
- Real-time updates (if WebSocket support is added)

## Project Structure

```
src/
  App.tsx              # Main app component with Refine setup
  dataProvider.ts      # Gogo data provider
  authProvider.ts      # Authentication provider (if needed)
  pages/
    users/
      list.tsx         # User list page
      create.tsx       # Create user page
      edit.tsx         # Edit user page
      show.tsx         # Show user page
    articles/
      ...
```

## Customization

You can customize the UI by:

1. Using different Refine UI frameworks (Ant Design, Material UI, Chakra UI, etc.)
2. Customizing field renderers
3. Adding custom actions
4. Implementing custom filters

See [Refine documentation](https://refine.dev/docs) for more details.

