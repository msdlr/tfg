pragma solidity >=0.4.22 <0.7.0;
import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/GeneralContract.sol";
import "./contractCreator.sol";

contract userFunctions_test{

    contractCreator cc;
    GeneralContract testContract;
    AuthContract ac;
    address thisContract = address(this);

    uint16 pass;

    function test_retrieveContract() public {
        cc = new contractCreator();
        testContract = cc.createContract();
        cc.addUser();
        ac = testContract.getUserAuthContract(thisContract);
    }

    function test_retrieved() public {
        // Check that this contract is a plain user in the new contract
        Assert.isTrue(testContract.getUserRegistered(thisContract),"[test_retrieved] User entry (isRegistered) should be TRUE");
        Assert.isFalse(testContract.getUserAdmin(thisContract),"[test_retrieved] User entry (adminStatus) should be FALSE");
        Assert.isFalse(testContract.amILocked(),"[test_locked] user is locked");
    }

    function test_getOTP() public {
        pass = testContract.getOTP();
        Assert.isTrue(pass != uint16(0),"[test_getOTP] Pass not retrieved");
    }

    //function test_trylogin() public {
    //    //testContract.tryLogin(pass);
    //    ac.tryLogin(pass);
    //    
    //    //Assert.isTrue(testContract.getUserLoggedIn(thisContract),"[test_trylogin] user was not logged in");
    //    //Assert.isTrue(testContract.getUserAttempts(thisContract) == 0,"[test_trylogin] user was not logged in");
    //}
}