describe('Login form', () => {
  it('allows a user to fill in login credentials', () => {
    cy.visit('/login');

    cy.contains('h1', 'Login').should('be.visible');
    cy.get('#login-email').type('samarth@example.com');
    cy.get('#login-password').type('password123');

    cy.get('#login-email').should('have.value', 'samarth@example.com');
    cy.get('#login-password').should('have.value', 'password123');
    cy.contains('button', 'Sign in').should('not.be.disabled');
  });
});
