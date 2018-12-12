/*
* Copyright 2018 PoC-Consortium
*/

import { Injectable } from "@angular/core";
import { RequestOptions, Response } from "@angular/http";
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import 'rxjs/add/operator/toPromise';
import 'rxjs/add/operator/timeout'

import { Account, Currency, EncryptedMessage, HttpError, Keys, Message, Settings, Transaction, constants } from "../model";
import { NoConnectionError, UnknownAccountError } from "../model/error";
import { BurstUtil } from "../util"
import { CryptoService } from "./crypto.service";
import { NotificationService} from "./notification.service";
import { StoreService } from "./store.service";
import { HttpClient, HttpHeaders, HttpParams } from "@angular/common/http";
import { Observable } from "rxjs";
import { fromPromise } from "rxjs/internal-compatibility";
import { tap } from "rxjs/operators";

/*
* AccountService class
*
* The AccountService is responsible for communication with the Burst node.
* It also preserves the current selected account and shares account related information
* across components.
*/
@Injectable()
export class AccountService {
    private nodeUrl: string;

    // Behaviour Subject for the current selected account, can be subscribed by components
    public currentAccount: BehaviorSubject<any> = new BehaviorSubject(undefined);

    constructor(
        private http: HttpClient,
        private cryptoService: CryptoService,
        private storeService: StoreService
    ) {
        this.storeService.settings.subscribe((settings: Settings) => {
            this.nodeUrl = settings.node;
        });
        // this.nodeUrl = "https://wallet.burst.cryptoguru.org:8125/burst"
    }

    public setCurrentAccount(account: Account) {
        this.currentAccount.next(account);
    }

    /*
    * Method responsible for creating a new active account from a passphrase.
    * Generates keys for an account, encrypts them with the provided key and saves them.
    * TODO: error handling of asynchronous method calls
    */
    public createActiveAccount({passphrase, pin = ""}): Promise<Account> {
        return new Promise((resolve, reject) => {
            let account: Account = new Account();
            // import active account
            account.type = "active";
            return this.cryptoService.generateMasterKeys(passphrase)
                .then(keys => {
                    let newKeys = new Keys();
                    newKeys.publicKey = keys.publicKey;
                    return this.cryptoService.encryptAES(keys.signPrivateKey, this.hashPinEncryption(pin))
                        .then(encryptedKey => {
                            newKeys.signPrivateKey = encryptedKey;
                            return this.cryptoService.encryptAES(keys.agreementPrivateKey, this.hashPinEncryption(pin))
                                .then(encryptedKey => {
                                    newKeys.agreementPrivateKey = encryptedKey;
                                    account.keys = newKeys;
                                    account.pinHash = this.hashPinStorage(pin, keys.publicKey);
                                    return this.cryptoService.getAccountIdFromPublicKey(keys.publicKey)
                                        .then(id => {
                                            account.id = id;
                                            return this.cryptoService.getBurstAddressFromAccountId(id)
                                                .then(address => {
                                                    account.address = address;
                                                    return this.storeService.saveAccount(account)
                                                        .then(account => {
                                                            resolve(account);
                                                        });
                                                });
                                        });
                                });
                        });
                });
        });
    }

    /*
    * Method responsible for importing an offline account.
    * Creates an account object with no keys attached.
    */
    public createOfflineAccount(address: string): Promise<Account> {
        return new Promise((resolve, reject) => {

            if (!BurstUtil.isValid(address)) {
                reject("Invalid Burst Address");
            }

            let account: Account = new Account();
            
            this.storeService.findAccount(BurstUtil.decode(address))
                .then(found => {
                    if (found == undefined) {
                        // import offline account
                        account.type = "offline";
                        account.address = address;
                        return this.cryptoService.getAccountIdFromBurstAddress(address)
                            .then(id => {
                                account.id = id;
                                return this.storeService.saveAccount(account)
                                    .then(account => {
                                        resolve(account);
                                    });
                            });
                    } else {
                        reject("Burstcoin address already imported!");
                    }
                })
        });
    }

    /*
    * Method responsible for activating an offline account.
    * This method adds keys to an existing account object and enables it.
    */
    public activateAccount(account: Account, passphrase: string, pin: string): Promise<Account> {
        return new Promise((resolve, reject) => {
            this.cryptoService.generateMasterKeys(passphrase)
                .then(keys => {
                    let newKeys = new Keys();
                    newKeys.publicKey = keys.publicKey;
                    return this.cryptoService.encryptAES(keys.signPrivateKey, this.hashPinEncryption(pin))
                        .then(encryptedKey => {
                            newKeys.signPrivateKey = encryptedKey;
                            return this.cryptoService.encryptAES(keys.agreementPrivateKey, this.hashPinEncryption(pin))
                                .then(encryptedKey => {
                                    newKeys.agreementPrivateKey = encryptedKey;
                                    account.keys = newKeys;
                                    account.pinHash = this.hashPinStorage(pin, keys.publicKey);
                                    account.type = "active";
                                    return this.storeService.saveAccount(account)
                                        .then(account => {
                                            resolve(account);
                                        });
                                });
                        });
                })
        });
    }

    /*
    * Method responsible for removing an existing account.
    */
    public removeAccount(account: Account): Promise<boolean> {
        return new Promise((resolve, reject) => {
            this.storeService.removeAccount(account)
                .then(success => {
                    resolve(success);
                })
                .catch(error => {
                    reject(error);
                })
        });
    }

    /*
    * Method responsible for synchronizing an account with the blockchain.
    */
    public synchronizeAccount(account: Account): Promise<Account> {
        return new Promise((resolve, reject) => {
            this.getBalance(account.id)
                .then(balance => {
                    account.balance = balance.confirmed;
                    account.unconfirmedBalance = balance.unconfirmed;
                    this.getTransactions(account.id)
                        .then(transactions => {
                            account.transactions = transactions;
                            this.getUnconfirmedTransactions(account.id)
                                .then(transactions => {
                                    account.transactions = transactions.concat(account.transactions);
                                    this.storeService.saveAccount(account)
                                        .catch(error => { console.log("Failed saving the account!"); })
                                    resolve(account);
                                }).catch(error => reject(error))
                        }).catch(error => reject(error))
                }).catch(error => reject(error))
        });
    }

    /*
    * Method responsible for selecting a different account.
    */
    public selectAccount(account: Account): Promise<Account> {
        return new Promise((resolve, reject) => {
            this.storeService.selectAccount(account)
                .then(account => {
                    this.synchronizeAccount(account);
                })
            this.setCurrentAccount(account);
            resolve(account);
        });
    }

    /*
    * Method responsible for getting the latest 15 transactions.
    */
    public getTransactions(id: string): Promise<Transaction[]> {
        return new Promise((resolve, reject) => {
            let params: HttpParams = new HttpParams()
                .set("requestType", "getAccountTransactions")
                .set("firstIndex", "0")
                .set("lastIndex", constants.transactionCount)
                .set("account", id);

            let requestOptions = BurstUtil.getRequestOptions();
            requestOptions.params = params;
            return this.http.get(this.nodeUrl, requestOptions)
                .timeout(constants.connectionTimeout)
                .toPromise<any>() // todo
                .then(response => {
                    if (response.errorCode) {
                        return reject(response.errorDescription);
                    }
                    let transactions: Transaction[] = [];
                    response.transactions.map(transaction => {
                        transaction.amountNQT = BurstUtil.convertStringToNumber(transaction.amountNQT);
                        transaction.feeNQT = BurstUtil.convertStringToNumber(transaction.feeNQT);
                        transactions.push(new Transaction(transaction));
                    });
                    return resolve(transactions);
                })
                .catch(error => {console.log(error);reject(new NoConnectionError())});
        });
    }

    /*
    * Method responsible for getting yet unconfirmed transactions.
    */
    public getUnconfirmedTransactions(id: string): Promise<Transaction[]> {
        return new Promise((resolve, reject) => {
            let params: HttpParams = new HttpParams()
                .set("requestType", "getUnconfirmedTransactions")
                .set("account", id);
            let requestOptions = BurstUtil.getRequestOptions();
            requestOptions.params = params;
            return this.http.get(this.nodeUrl, requestOptions)
                .timeout(constants.connectionTimeout)
                .toPromise<any>() // todo
                .then(response => {
                    let transactions: Transaction[] = [];
                    response.unconfirmedTransactions.map(transaction => {
                        transaction.amountNQT = BurstUtil.convertStringToNumber(transaction.amountNQT);
                        transaction.feeNQT = BurstUtil.convertStringToNumber(transaction.feeNQT);
                        transaction.confirmed = false;
                        transactions.push(new Transaction(transaction));
                    });
                    resolve(transactions);
                })
                .catch(error => reject(new NoConnectionError()));
        });
    }

    /*
    * Method responsible for getting one specific transaction
    */
    public getTransaction(id: string): Promise<Transaction> {
        return new Promise((resolve, reject) => {
            let params: HttpParams = new HttpParams()
                .set("requestType", "getTransaction")
                .set("transaction", id);
            let requestOptions = BurstUtil.getRequestOptions();
            requestOptions.params = params;
            return this.http.get(this.nodeUrl, requestOptions)
                .timeout(constants.connectionTimeout)
                .toPromise<any>() // todo
                .then(response => {
                    return response || [];
                })
                .catch(error => reject(new NoConnectionError()));
        });
    }

    /*
    * Method responsible for getting the current balance of an account.
    */
    public getBalance(id: string): Promise<any> {
        return new Promise((resolve, reject) => {
            let params: HttpParams = new HttpParams()
                .set("requestType", "getBalance")
                .set("account", id);
            let requestOptions = BurstUtil.getRequestOptions();
            requestOptions.params = params;
            return this.http.get(this.nodeUrl, requestOptions)
                .timeout(constants.connectionTimeout)
                .toPromise<any>() // todo
                .then(response => {
                    if (response.errorCode == undefined) {
                        let balanceString = response.guaranteedBalanceNQT;
                        balanceString = BurstUtil.convertStringToNumber(balanceString);
                        let unconfirmedBalanceString = response.unconfirmedBalanceNQT;
                        unconfirmedBalanceString = BurstUtil.convertStringToNumber(unconfirmedBalanceString);
                        resolve({ confirmed: parseFloat(balanceString), unconfirmed: parseFloat(unconfirmedBalanceString) });
                    } else {
                        if (response.errorDescription == "Unknown account") {
                            reject(new UnknownAccountError())
                        } else {
                            reject(new Error("Failed fetching balance"));
                        }
                    }
                })
                .catch(error => reject(new NoConnectionError()));
        });
    }

    /*
    * Method responsible for getting the public key in the blockchain of an account.
    */
    public getAccountPublicKey(id: string): Promise<string> {
        return new Promise((resolve, reject) => {
            let params: HttpParams = new HttpParams()
                .set("requestType", "getAccountPublicKey")
                .set("account", id);
            let requestOptions = BurstUtil.getRequestOptions();
            requestOptions.params = params;
            return this.http.get(this.nodeUrl, requestOptions)
                .timeout(constants.connectionTimeout)
                .toPromise<any>() // todo
                .then(response => {
                    if (response.publicKey != undefined) {
                        let publicKey = response.publicKey;
                        resolve(response.publicKey);
                    } else {
                        reject(new UnknownAccountError())
                    }
                })
                .catch(error => reject(new NoConnectionError()));
        });
    }

    /*
    * Method responsible for executing a transaction.
    */
    public doTransaction(transaction: Transaction, encryptedPrivateKey: string, pin: string): Promise<Transaction> {
        return new Promise((resolve, reject) => {
            let params: HttpParams = new HttpParams()
                .set("requestType", "sendMoney")
                .set("amountNQT", BurstUtil.convertNumberToString(transaction.amountNQT))
                .set("deadline", "1440") // todo
                .set("feeNQT", BurstUtil.convertNumberToString(transaction.feeNQT))
                .set("publicKey", transaction.senderPublicKey)
                .set("recipient", transaction.recipientAddress);
            if (transaction.attachment != undefined) {
                params = this.constructAttachment(transaction, params);
            }
            let requestOptions = BurstUtil.getRequestOptions();
            requestOptions.params = params;
            // request 'sendMoney' to burst node
            return this.http.post(this.nodeUrl, {}, requestOptions)
                .timeout(constants.connectionTimeout)
                .toPromise<any>() // todo
                .then(this.postTransaction(resolve, reject, transaction, encryptedPrivateKey, pin))
                .catch(error => { console.log(error); reject("Transaction error: Generating transaction. Check the recipient!") });
        });
    }


    /*
    * Method responsible for sending a message
    */
   public sendMessage(transaction: Transaction, encryptedPrivateKey: string, pin: string): Promise<Transaction> {
        return new Promise((resolve, reject) => {

            let params: HttpParams = new HttpParams()
                .set("requestType", "sendMessage")
                .set("deadline", "1440") // todo
                .set("feeNQT", BurstUtil.convertNumberToString(transaction.feeNQT))
                .set("publicKey", transaction.senderPublicKey)
                .set("recipient", transaction.recipientAddress);

            if (transaction.attachment != undefined) {
                params = this.constructAttachment(transaction, params);
            }

            let requestOptions = BurstUtil.getRequestOptions();
            requestOptions.params = params;

            return this.http.post(this.nodeUrl, {}, requestOptions)
                .timeout(constants.connectionTimeout)
                .toPromise<any>() // todo
                .then(this.postTransaction(resolve, reject, transaction, encryptedPrivateKey, pin))
                .catch(error => { console.log(error); reject("Transaction error: Generating transaction. Check the recipient!") });
        });
    }

    private postTransaction: any = (resolve, reject, transaction, encryptedPrivateKey, pin) => async (response) => {
        if (response.unsignedTransactionBytes != undefined) {
            // get unsigned transactionbytes
            const unsignedTransactionHex = response.unsignedTransactionBytes;
            // sign unsigned transaction bytes
            const signature = await this.cryptoService.generateSignature(unsignedTransactionHex, encryptedPrivateKey, this.hashPinEncryption(pin));
            const verified = await this.cryptoService.verifySignature(signature, unsignedTransactionHex, transaction.senderPublicKey);
            if (verified) {
                const signedTransactionBytes = await this.cryptoService.generateSignedTransactionBytes(unsignedTransactionHex, signature);
                const params = new HttpParams()
                    .set("requestType", "broadcastTransaction")
                    .set("transactionBytes", signedTransactionBytes);
                let requestOptions = BurstUtil.getRequestOptions();
                requestOptions.params = params;
                // request 'broadcastTransaction' to burst node
                return this.http.post(this.nodeUrl, {}, requestOptions)
                    .timeout(constants.connectionTimeout)
                    .toPromise<any>() // todo
                    .then(response => {
                        const params = new HttpParams()
                            .set("requestType", "getTransaction")
                            .set("transaction", response.transaction);
                        requestOptions = BurstUtil.getRequestOptions();
                        requestOptions.params = params;
                        // request 'getTransaction' to burst node
                        return this.http.get(this.nodeUrl, requestOptions)
                            .timeout(constants.connectionTimeout)
                            .toPromise<any>() // todo
                            .then(response => {
                                resolve(new Transaction(response));
                            })
                            .catch(error => reject("Transaction error: Finalizing transaction!"));
                    }).catch(error => reject("Transaction error: Executing transaction!"));
            }
            else {
                reject("Transaction error: Verifying signature!");
            }
        }
        else {
            reject("Transaction error: Generating transaction. Check the recipient!");
        }
    };

    private constructAttachment(transaction: Transaction, params: HttpParams) {
        if (transaction.attachment.type == "encrypted_message") {
            let em: EncryptedMessage = <EncryptedMessage>transaction.attachment;
            params = params.set("encryptedMessageData", em.data)
                .set("encryptedMessageNonce", em.nonce)
                .set("messageToEncryptIsText", String(em.isText));
        }
        else if (transaction.attachment.type == "message") {
            let m: Message = <Message>transaction.attachment;
            params = params.set("message", m.message)
                .set("messageIsText", String(m.messageIsText));
        }
        return params;
    }

    public doMultiOutTransaction(transaction: Transaction, encryptedPrivateKey: string, pin: string, sameAmount: boolean): Promise<Transaction> {
        return new Promise((resolve, reject) => {

            let params: HttpParams = new HttpParams()
                .set("requestType", sameAmount ? "sendMoneyMultiSame" : "sendMoneyMulti")
                .set('feeNQT', transaction.feeNQT.toString())
                .set('deadline', transaction.deadline)
                .set('recipients', transaction.recipients)
                .set('publicKey', transaction.senderPublicKey);

            if (sameAmount) {
                params = params.set('amountNQT', transaction.amountNQT.toString());
            }

            let requestOptions = BurstUtil.getRequestOptions();
            requestOptions.params = params;

            return this.http.post(this.nodeUrl, {}, requestOptions)
                .timeout(constants.connectionTimeout)
                .toPromise<any>() // todo
                .then(this.postTransaction(resolve, reject, transaction, encryptedPrivateKey, pin))
                .catch(error => { console.log(error); reject("Transaction error: Generating transaction. Check the recipient!") });
        });
    }

    handleIncomingTransactions(transactions: Transaction[]) {
        console.log(`NEW TRANSACTIONS FOUND`, transactions); // todo
    }

    /*
    * Method responsible for verifying the PIN
    */
    public checkPin(pin: string): boolean {
        return this.currentAccount.value != undefined ? this.currentAccount.value.pinHash == this.hashPinStorage(pin, this.currentAccount.value.keys.publicKey) : false;
    }

    /*
    * Method responsible for hashing the PIN to carry out an ecryption.
    */
    public hashPinEncryption(pin: string): string {
        // todo: make this work?
        // if (this.currentAccount.value) {
        //     pin = pin + this.currentAccount.value.id;
        // }
        return this.cryptoService.hashSHA256(pin);
    }

    /*
    * Method responsible for hashing the PIN for saving it into the database.
    */
    public hashPinStorage(pin: string, publicKey: string): string {
        return this.cryptoService.hashSHA256(pin + publicKey);
    }

    /*
    * Method responsible for pin validation.
    */
    public isPin(pin: string): boolean {
        return /^[0-9]{6}$/i.test(pin);
    }
}
