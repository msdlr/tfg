pragma solidity >=0.4.22 <0.7.0;
import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/GeneralContract.sol";

// File name has to end with '_test.sol', this file can contain more than one testSuite contracts
contract General_test {
    
    GeneralContract testContract;
    address thisContract = address(this);
    address me = msg.sender;
    
    address[] testAddrs = 
    [0xAb8483F64d9C6d1EcF9b849Ae677dD3315835cb2,
    0x4B20993Bc481177ec7E8f571ceCaE8A9e22C02db,
    0x78731D3Ca6b7E34aC0F824c42a7cC18A495cabaB,
    0x617F2E2fD72FD9D5503197092aC168c91465E7f2,
    0x17F6AD8Ef982297579C203069C1DbfFE4348c372,
    0x5c6B0f7Bf3E7ce046039Bd8FABdfD3f9F5021678,
    0x03C6FcED478cBbC9a4FAB34eF9f40767739D1Ff7,
    0x03C6FcED478cBbC9a4FAB34eF9f40767739D1Ff7,
    0x0A098Eda01Ce92ff4A4CCb7A4fFFb5A43EBC70DC,
    0xCA35b7d915458EF540aDe6068dFe2F44E8fa733c,
    0x14723A09ACff6D2A60DcdF7aA4AFf308FDDC160C,
    0x4B0897b0513fdC7C541B6d9D7E929C4e5364D2dB,
    0x583031D1113aD414F02576BD6afaBfb302140225,
    0xdD870fA1b7C4700F2BD7f44238821C26f7392148];
    

    /// 'beforeAll' runs before all other tests
    /// More special functions are: 'beforeEach', 'beforeAll', 'afterEach' & 'afterAll'

    // msg.sender in GeneralContract -> address (this)
    function test_createContract() public {
        // Instantiate the contract to test
        testContract = new GeneralContract(thisContract, "11223344K");
        // Check that the constructor executed correctly
        Assert.equal(thisContract, testContract.getOwner(), "owner address should be this caller");
        
        // Check if the fields initialized correctly
        Assert.isTrue(testContract.getUserAdmin(thisContract),"caller was not made admin");
        Assert.isTrue(testContract.getUserRegistered(thisContract),"caller is not on the user list");
        Assert.isFalse(testContract.getUserLoggedIn(thisContract),"caller is logged in");
        Assert.equal("11223344K", testContract.getUserId(thisContract), "contract id was not set");
    }

    function test_getters() public {
        Assert.equal(address(testContract), testContract.getContractAddress(), "contract address differs");
        Assert.equal("11223344K", testContract.getUserId(thisContract), "id differs");
        Assert.isTrue(testContract.getUserRegistered(thisContract), "caller not registered");
        Assert.isFalse(testContract.getUserLoggedIn(thisContract), "caller should not be logged in");
        Assert.isTrue(testContract.getUserAdmin(thisContract), "caller should be admin");
        Assert.equal(uint(0), uint(testContract.getUserAttempts(thisContract)), "attempts should be 0");
        Assert.isFalse(testContract.amILocked(), "should return false");
    }

    function test_addUser() public {
        
        // Try to add a new user
        Assert.isFalse(testContract.getUserRegistered(testAddrs[0]),"User entry (isRegistered) should be FALSE");
        testContract.addUser(testAddrs[0], "11223344K");
        
        // Check it was done correctly
        Assert.isTrue(testContract.getUserRegistered(testAddrs[0]),"User entry (isRegistered) should be TRUE");
    }

    function test_promoteUser() public {
        Assert.isFalse(testContract.getUserAdmin(testAddrs[0]),"User entry (adminStatus) should be FALSE");
        testContract.promoteUser(testAddrs[0]);
        Assert.isTrue(testContract.getUserAdmin(testAddrs[0]),"User entry (adminStatus) should be FALSE");
    }

    function test_demoteUser() public {
        Assert.isTrue(testContract.getUserAdmin(testAddrs[0]),"User entry (adminStatus) should be FALSE");
        testContract.demoteAdmin(testAddrs[0]);
        Assert.isFalse(testContract.getUserAdmin(testAddrs[0]),"User entry (adminStatus) should be FALSE");
    }

    
}