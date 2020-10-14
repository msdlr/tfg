pragma solidity >=0.6.4 <=7.3.0;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/GeneralContract.sol";

// File name has to end with '_test.sol', this file can contain more than one testSuite contracts
contract login_test {
    
    GeneralContract testContract;
    address thisContract = address(this);
    uint16 pass;

    /// 'beforeAll' runs before all other tests
    /// More special functions are: 'beforeEach', 'beforeAll', 'afterEach' & 'afterAll'

    // msg.sender in GeneralContract -> address (this) in this contract
    function test_createContract() public {
        // Instantiate the contract to test
        testContract = new GeneralContract(thisContract, "M4573R");
        // Check that the constructor executed correctly
        //Assert.equal(thisContract, testContract.getOwner(), "owner address should be this caller");
        
        // Check if the fields initialized correctly
        Assert.isTrue(testContract.getUserAdmin(thisContract),"[test_createContract] caller was not made admin");
        Assert.isTrue(testContract.getUserRegistered(thisContract),"[test_createContract] caller is not on the user list");
        Assert.isFalse(testContract.getUserLoggedIn(thisContract),"[test_createContract] caller is logged in");
        Assert.equal("M4573R", testContract.getUserId(thisContract), "[test_createContract] contract id was not set");
    }

    function test_getOTP() public {
        pass = testContract.getOTP();
        Assert.isTrue(pass != uint16(0),"[test_getOTP] Pass not retrieved");
    }

    function test_trylogin() public {
        testContract.tryLogin(pass);
        Assert.isTrue(testContract.getUserLoggedIn(thisContract),"[test_trylogin] user was not logged in");
        Assert.isTrue(testContract.getUserAttempts(thisContract) == 0,"[test_trylogin] user was not logged in");
    }

    function test_trylogout() public {
        testContract.tryLogout();
        Assert.isFalse(testContract.getUserLoggedIn(thisContract),"[test_trylogout] logout failed");
    }
} 
