# coinage

Golang toolkit for a text-based double-entry accounting system

### Balancing method

- transaction has multiple "lines"
- each line has an account reference, a "value", and an optional "price"
- the price allows a transaction to be balanced when different commodities are used    

    ```
    # example of converting Pounds to Rand using cash
    
    Assets:Cash   183.25 ZAR for 10 GBP
    Assets:Cash    -10 GBP
    ```
